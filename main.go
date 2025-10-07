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
	// 1) Load env duluan
	config.LoadEnv()

	// 2) Connect DB
	database.ConnectDB()
	defer database.DB.Close()

	// 3) Fiber app + middleware dasar
	app := fiber.New(fiber.Config{
		// optional: batas body biar aman
		BodyLimit: 10 * 1024 * 1024, // 10MB
		// optional: custom error handler sederhana
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

	// Logger hanya saat non-production
	if os.Getenv("APP_ENV") != "production" {
		app.Use(logger.New())
	}
	app.Use(recover.New())
	app.Use(cors.New())

	// 4) (Opsional) Root endpoint singkat — atau pindahkan ke router, pilih salah satu
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Alumni API")
	})

	// 5) Register routes
	route.SetupRoutes(app, database.DB)

	// 6) Start server + graceful shutdown
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Jalankan server di goroutine
	go func() {
		log.Printf("Server listening on :%s\n", port)
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Server stopped: %v\n", err)
		}
	}()

	// Tangkap signal OS
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Shutdown dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}
}
