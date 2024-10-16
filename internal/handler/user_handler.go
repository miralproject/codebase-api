package handler

import (
	"codebase-api/config"
	"codebase-api/internal/domain"
	"codebase-api/internal/usecase"
	helper "codebase-api/pkg/helpers"
	middleware "codebase-api/pkg/middlewares"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cast"
)

type UserResponseDto struct {
	ID        int    `json:"id"`
	UUID      string `json:"uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	IsActive  bool   `json:"is_active"`
}

func ToUserResponseDto(user interface{}) UserResponseDto {
	var u domain.User

	// Cek apakah user adalah pointer atau nilai
	switch v := user.(type) {
	case *domain.User:
		u = *v // Dereferensiasi pointer
	case domain.User:
		u = v // Langsung gunakan nilai
	}

	return UserResponseDto{
		ID:        int(u.ID),
		UUID:      u.UUID.String(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		IsActive:  u.IsActive,
	}
}

type UserHandler struct {
	usecase  *usecase.UserUseCase
	validate *validator.Validate
}

func NewUserHandler(usecase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{usecase: usecase, validate: validator.New()}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
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
		Password:  input.Password, // Harus di-hash di use case
		Phone:     input.Phone,
	}

	err := h.usecase.Register(user)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Internal servel error", nil)
	}

	return helper.SuccessResponse(c, nil, "Register successful")
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
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
		fmt.Println(err)
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

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return helper.SuccessResponse(c, nil, "Logout successful")
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
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

func (h *UserHandler) All(c *fiber.Ctx) error {
	users, err := h.usecase.FinAll()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Failed to fetch users", err)

	}
	var dto []UserResponseDto
	for _, user := range users {
		dto = append(dto, ToUserResponseDto(user))
	}

	return helper.SuccessResponse(c, dto, "Fetch all data users success")
}

func (h *UserHandler) Searching(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := c.Query("search", "")
	isActiveParam := c.Query("status", "")

	var isActive *bool
	if isActiveParam != "" {
		active := isActiveParam == "true"
		isActive = &active
	}

	users, err := h.usecase.Searching(isActive, search)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Failed to fetch users", err)
	}

	totalItems := len(users)
	totalPages := (totalItems + pageSize - 1) / pageSize

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > totalItems {
		start = totalItems
	}
	if end > totalItems {
		end = totalItems
	}

	paginatedUsers := users[start:end]

	var dto []UserResponseDto
	for _, user := range paginatedUsers {
		dto = append(dto, ToUserResponseDto(user))
	}

	var data []interface{}
	for _, user := range dto {
		data = append(data, user)
	}

	return helper.PaginationResponse(c, data, page, pageSize, totalItems, totalPages, "Fetch all data users success")
}

func (h *UserHandler) Detail(c *fiber.Ctx) error {
	user, err := h.usecase.FindById(cast.ToUint(c.Params("id")))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, "User not found", err)
	}

	dto := ToUserResponseDto(user)
	return helper.SuccessResponse(c, dto, "Fetch data users success")
}
