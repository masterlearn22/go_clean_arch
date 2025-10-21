package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_clean/config"
	"go_clean/database"
	"go_clean/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1️⃣ Load env
	config.LoadEnv()

	// 2️⃣ Connect ke PostgreSQL
	database.ConnectDB()
	defer database.DB.Close()

	// 3️⃣ Connect ke MongoDB
	database.ConnectMongoDB()

	// 4️⃣ Setup Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// 5️⃣ Middleware
	if os.Getenv("APP_ENV") != "production" {
		app.Use(logger.New())
	}
	app.Use(recover.New())
	app.Use(cors.New())

	// 6️⃣ Root sederhana
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Alumni API 🚀")
	})

	// 7️⃣ Register routes (Postgres + Mongo)
	route.SetupPekerjaanMongoRoutes(app, database.MongoDB)
	route.SetupAlumniMongoRoutes(app, database.MongoDB)
	route.SetupRoutes(app, database.DB, database.MongoDB)


	// 8️⃣ Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		log.Printf("Server running on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	// 9️⃣ Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}
