package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Initialize logrus logger
var lg = logrus.New()

func SetupLogger() {
	lg.SetFormatter(&logrus.JSONFormatter{})
	lg.SetLevel(logrus.InfoLevel)
}

// LogrusMiddleware logs every HTTP request
func LogrusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Logging HTTP requests
		lg.WithFields(logrus.Fields{
			"method": c.Method(),
			"path":   c.Path(),
			"status": c.Response().StatusCode(),
			"took":   time.Since(start),
		}).Info("request handled")

		return err
	}
}
