package route

import (
	"context"
	"time"

	"go_clean/app/models"
	"go_clean/app/repository"
	"go_clean/app/service"
	"go_clean/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupAlumniMongoRoutes(app *fiber.App, mongoDB *mongo.Database) {
	// 🔧 Inisialisasi repository & service
	repo := repository.NewAlumniMongoRepository(mongoDB)
	svc := service.NewAlumniMongoService(repo)

	// 🧩 Semua endpoint butuh login (AuthRequired)
	api := app.Group("/api/alumni-mongo", middleware.AuthRequired())

	// ========== READ (bisa diakses semua user login) ==========

	// GET /api/alumni-mongo → Ambil semua data alumni
	api.Get("/", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := svc.GetAll(ctx)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	// GET /api/alumni-mongo/:id → Ambil 1 data alumni by ID
	api.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := svc.GetByID(ctx, id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	// ========== ADMIN ONLY (CREATE, UPDATE, DELETE) ==========

	admin := api.Group("", middleware.AdminOnly())

	// POST /api/alumni-mongo → Tambah data (hanya admin)
	admin.Post("/", func(c *fiber.Ctx) error {
		var input models.AlumniMongo
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "JSON tidak valid"})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := svc.Create(ctx, &input)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(201).JSON(data)
	})

	// PUT /api/alumni-mongo/:id → Update data (hanya admin)
	admin.Put("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var input models.AlumniMongo

		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "JSON tidak valid"})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := svc.Update(ctx, id, &input)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(data)
	})

	// DELETE /api/alumni-mongo/:id → Hapus data (hanya admin)
	admin.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := svc.Delete(ctx, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Data alumni berhasil dihapus"})
	})
}
