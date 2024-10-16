package router

import (
	"codebase-api/config/rabbitmq"
	"codebase-api/internal/handler"
	"codebase-api/internal/repository"
	"codebase-api/internal/usecase"
	middleware "codebase-api/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, ch *amqp.Channel) {
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	api := app.Group("/api/v1")
	api.Get("/healty", func(c *fiber.Ctx) error { return c.SendString("healty is good!!") })

	api.Post("/auth/register", userHandler.Register)
	api.Post("/auth/login", userHandler.Login)
	api.Post("/auth/logout", userHandler.Logout)
	api.Post("/auth/refresh-token", middleware.JwtProtected(), userHandler.RefreshToken)

	api.Get("/users", middleware.JwtProtected(), userHandler.All)
	api.Get("/users/search", middleware.JwtProtected(), userHandler.Searching)
	api.Get("/users/:id", middleware.JwtProtected(), userHandler.Detail)

	api.Post("/publish", func(c *fiber.Ctx) error {
		type RequestBody struct {
			Message string `json:"message"`
		}

		var body RequestBody
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		err := rabbitmq.PublishMessage(ch, "testQueue", body.Message)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to publish message")
		}

		return c.SendString("Message published to RabbitMQ")
	})
}
