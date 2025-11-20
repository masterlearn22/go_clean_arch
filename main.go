package main

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @tag.name Auth
// @tag.description Endpoint untuk autentikasi dan login user
// @tag.order 1

// @tag.name Admin
// @tag.description Endpoint untuk Membuat user baru (Admin Only)
// @tag.order 2

// @tag.name Alumni-Mongo
// @tag.description Endpoint untuk data alumni
// @tag.order 3

// @tag.name Pekerjaan-Mongo
// @tag.description Endpoint untuk data pekerjaan alumni (MongoDB)
// @tag.order 4

// @tag.name FileUpload
// @tag.description Upload, lihat dan hapus file
// @tag.order 5

// @tag.name Alumni-PostgresSQL
// @tag.description Upload, lihat dan hapus file
// @tag.order 6

// @tag.name Pekerjaan-PostgresSQL
// @tag.description Endpoint untuk data pekerjaan alumni (PostgreSQL)
// @tag.order 7	

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	// "fmt"
	// "go_clean/utils"

	"go_clean/config"
	"go_clean/database"
	FiberApp "go_clean/fiber"
	routePostgre "go_clean/route/postgresql"
	routeMongo "go_clean/route/mongodb"

	// "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	 fiberSwagger "github.com/swaggo/fiber-swagger"

	docs "go_clean/docs"
	// _ "go_clean/docs"
)

// @title User API
// @version 1.0
// @description API untuk mengelola data user dengan MongoDB menggunakan Clean Architecture
// @host localhost:3000
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization


func main() {
	// 1️ Load env
	config.LoadEnv()

	// 2️ Connect ke PostgreSQL
	database.ConnectDB()
	defer database.DB.Close()

	// 3️ Connect ke MongoDB
	database.ConnectMongoDB()

	// 4️ Setup Fiber app
	app := FiberApp.SetupFiber()

	// Swagger
	docs.SwaggerInfo.BasePath = "/api"
    app.Get("/swagger/*", fiberSwagger.WrapHandler)
	log.Println("➡️  Swagger UI available at: http://localhost:" + os.Getenv("APP_PORT") + "/swagger/index.html")


	// 5️ Middleware
	if os.Getenv("APP_ENV") != "production" {
		app.Use(logger.New())
	}
	app.Use(recover.New())
	app.Use(cors.New())

	// 6️ Register  routes MongoDB
	routeMongo.SetupAuthMongoRoutes(app, database.MongoDB)
	routeMongo.SetupPekerjaanMongoRoutes(app, database.MongoDB)
	routeMongo.SetupAlumniMongoRoutes(app, database.MongoDB)
	routeMongo.SetupFileRoutes(app, database.MongoDB)

	// 7️ Register routes PostgreSQL
	routePostgre.SetupRoutes(app, database.DB)
	

	// 8 Start server
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

	// 9 Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}
