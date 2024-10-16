package main

import (
	"codebase-api/config"
	"codebase-api/config/rabbitmq"
	middleware "codebase-api/pkg/middlewares"
	"codebase-api/router"
	"log"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize Fiber db
	db := config.InitDB()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "expense-api v1.0.1",
	})

	// Initialize connecting and channel RabbitMQ
	conn, ch, err := rabbitmq.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	// Setup logger
	middleware.SetupLogger()

	// Register middleware
	app.Use(middleware.LogrusMiddleware())

	/*
	* middleware: Helmet
	* description: Helmet middleware helps secure your apps by setting various HTTP headers.
	 */
	app.Use(helmet.New())

	// Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint:  "/live",
		ReadinessEndpoint: "/ready",
	}))

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	// Register routes
	router.SetupRoutes(app, db, ch)

	// example consume mq
	// go func() {
	// 	msgs, err := rabbitmq.ConsumeMessages(ch, "testQueue")
	// 	if err != nil {
	// 		log.Fatalf("Failed to consume messages: %v", err)
	// 	}

	// 	// Proses pesan yang diterima
	// 	for msg := range msgs {
	// 		log.Printf("Received a message: %s", msg.Body)
	// 	}
	// }()

	// Start server
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
