package helpers

import (
	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(c *fiber.Ctx, data interface{}, message string) error {
	if data == nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"statusCode": 201,
			"message":    message,
			"data":       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"statusCode": 200,
		"message":    message,
		"data":       data,
	})
}

// Success Response with Pagination
func PaginationResponse(c *fiber.Ctx, data []interface{}, page int, pageSize int, totalItems int, totalPages int, message string) error {
	if data == nil {
		data = []interface{}{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"statusCode": 200,
		"message":    message,
		"data": fiber.Map{
			"items": data,
			"pagination": fiber.Map{
				"page":       page,
				"pageSize":   pageSize,
				"totalItems": totalItems,
				"totalPages": totalPages,
			},
		},
	})
}

// Error Response
func ErrorResponse(c *fiber.Ctx, statusCode int, message interface{}, err error) error {
	var errorMessage string
	if statusCode == fiber.StatusUnauthorized {
		errorMessage = "Unauthorized access"
		if err != nil {
			errorMessage = err.Error()
		}
	} else {
		if err != nil {
			errorMessage = err.Error()
		} else {
			errorMessage = "An unexpected error occurred"
		}
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"statusCode": statusCode,
		"message":    message,
		"error":      errorMessage,
	})
}
