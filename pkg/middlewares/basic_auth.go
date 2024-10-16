package middleware

import (
	"codebase-api/pkg/helpers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

func BasicAuth() fiber.Handler {
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Username or password not set in environment variables")
	}

	return basicauth.New(basicauth.Config{
		Authorizer: func(user, pass string) bool {
			return user == username && pass == password
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return helpers.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized access", nil)
		},
	})
}
