package route

import (
	"context"
	"strconv"
	"time"

	"go_clean/app/models"
	"go_clean/app/repository"
	"go_clean/app/service"
	"go_clean/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupPekerjaanMongoRoutes(app *fiber.App, mongoDB *mongo.Database) {
	repo := repository.NewPekerjaanMongoRepository(mongoDB)
	svc := service.NewPekerjaanMongoService(repo)

	// Semua endpoint butuh login
	api := app.Group("/api/pekerjaan-mongo", middleware.AuthRequired())

	// ========== READ (semua user login bisa) ==========
	api.Get("/", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := svc.GetAll(ctx)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

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

	// ========== ADMIN ONLY (READ + WRITE) ==========
	admin := api.Group("", middleware.AdminOnly())

	// GET /api/pekerjaan-mongo/alumni/:alumni_id → Admin only
	admin.Get("/alumni/:alumni_id", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("alumni_id"))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := svc.GetByAlumniID(ctx, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	// POST → Tambah data (hanya admin)
	admin.Post("/", func(c *fiber.Ctx) error {
		var input models.PekerjaanMongo
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "JSON tidak valid"})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := svc.Create(ctx, &input)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(201).JSON(result)
	})

	// PUT → Update data (hanya admin)
	admin.Put("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var input models.PekerjaanMongo

		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "JSON tidak valid"})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := svc.Update(ctx, id, &input)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(result)
	})

	// DELETE → Hapus data (hanya admin)
	admin.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := svc.Delete(ctx, id); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Data berhasil dihapus"})
	})
}
