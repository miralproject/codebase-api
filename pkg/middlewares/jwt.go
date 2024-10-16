package middleware

import (
	"codebase-api/config"
	"codebase-api/internal/domain"
	helper "codebase-api/pkg/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(config.GetJWTSecret())

func JwtProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("jwt")
		if cookie == "" {
			return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		token, _ := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
			}
			return jwtSecret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("id", claims["id"])
			c.Locals("uuid", claims["uuid"])
			c.Locals("username", claims["username"])
			return c.Next()
		}

		return helper.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
}

func GenerateJWT(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID,
		"uuid":     user.UUID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
