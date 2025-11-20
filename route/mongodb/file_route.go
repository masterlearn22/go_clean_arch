package route

import (
    Repo "go_clean/app/repository/mongodb"
    Srv "go_clean/app/service/mongodb"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
)

func SetupFileRoutes(app *fiber.App, mongoDB *mongo.Database) {

    // Folder static untuk akses file
    app.Static("/uploads", "./uploads")

    api := app.Group("/api")
    fileGroup := api.Group("/files")

    repo := Repo.NewFileRepository(mongoDB)
    srv := Srv.NewFileService(repo, "./uploads")

    fileGroup.Post("/upload", srv.UploadFile)
    fileGroup.Get("/", srv.GetAllFiles)
    fileGroup.Get("/:id", srv.GetFileByID)
    fileGroup.Delete("/:id", srv.DeleteFile)
}

