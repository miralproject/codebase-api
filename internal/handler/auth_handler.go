package handler

import (
	"codebase-api/config"
	"codebase-api/internal/domain"
	"codebase-api/internal/usecase"
	helper "codebase-api/pkg/helpers"
	middleware "codebase-api/pkg/middlewares"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cast"
)

type AuthHandler struct {
	usecase  *usecase.UserUseCase
	validate *validator.Validate
}

func NewAuthHandler(usecase *usecase.UserUseCase) *AuthHandler {
	return &AuthHandler{usecase: usecase, validate: validator.New()}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input struct {
		FirstName       string `json:"first_name" validate:"required"`
		LastName        string `json:"last_name" validate:"required"`
		Username        string `json:"username" validate:"required"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required,min=8"`
		PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
		Phone           string `json:"phone" validate:"required,e164"`
	}

	if err := c.BodyParser(&input); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	if err := h.validate.Struct(&input); err != nil {
		errorFields := helper.ValidationErrorFormatter(err, input)
		return helper.ErrorResponse(c, fiber.StatusBadRequest, errorFields, nil)
	}

	user := &domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Username:  input.Username,
		Email:     input.Email,
		Password:  input.Password,
		Phone:     input.Phone,
	}

	err := h.usecase.Register(user)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Internal servel error", nil)
	}

	return helper.SuccessResponse(c, nil, "Register successful")
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request", err)
	}

	if err := h.validate.Struct(&input); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "username and password are required", nil)
	}

	user, err := h.usecase.Login(input.Username, input.Password)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	token, err := middleware.GenerateJWT(*user)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Could not login", nil)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	dto := UserResponseDto{
		ID:        int(user.ID),
		UUID:      user.UUID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		IsActive:  user.IsActive,
	}

	// Return JWT token (for simplicity, skipped token creation)
	return helper.SuccessResponse(c, dto, "Login successful")
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return helper.SuccessResponse(c, nil, "Logout successful")
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return config.GetJWTSecret(), nil
	})

	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["id"].(float64)

		user, err := h.usecase.FindById(cast.ToUint(userID))
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "User not found", nil)
		}

		newToken, err := middleware.GenerateJWT(*user)
		if err != nil {
			return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Could not refresh token", nil)
		}

		// Set new token to cookies
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    newToken,
			Expires:  time.Now().Add(time.Hour * 24),
			HTTPOnly: true,
		})

		return helper.SuccessResponse(c, nil, "Token refreshed successful")
	}

	return helper.SuccessResponse(c, nil, "Token refreshed successful")
}
