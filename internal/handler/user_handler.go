package handler

import (
	"codebase-api/internal/domain"
	"codebase-api/internal/usecase"
	helper "codebase-api/pkg/helpers"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
