package route

import (
	"database/sql"
	"prak4/app/repository"
	"prak4/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	// Initialize repositories
	alumniRepo := &repository.AlumniRepository{DB: db}
	pekerjaanRepo := &repository.PekerjaanRepository{DB: db}

	// Initialize services
	alumniService := &service.AlumniService{Repo: alumniRepo}
	pekerjaanService := &service.PekerjaanService{Repo: pekerjaanRepo}


	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Alumni Management API")
	})

	// Grouping routes
	api := app.Group("/api")

	// Alumni Routes [cite: 588]
	alumni := api.Group("/alumni")
	alumni.Get("/", alumniService.GetAllAlumni)
	alumni.Get("/:id", alumniService.GetAlumniByID)
	alumni.Post("/", alumniService.CreateAlumni)
	alumni.Put("/:id", alumniService.UpdateAlumni)
	alumni.Delete("/:id", alumniService.DeleteAlumni)

	alumni.Get("/angkatan/:angkatan", alumniService.GetAlumniByAngkatan)

	// Pekerjaan Alumni Routes [cite: 594]
	pekerjaan := api.Group("/pekerjaan")
	pekerjaan.Get("/", pekerjaanService.GetAllPekerjaan)
	pekerjaan.Get("/:id", pekerjaanService.GetPekerjaanByID)
	pekerjaan.Get("/alumni/:alumni_id", pekerjaanService.GetPekerjaanByAlumniID)
	pekerjaan.Post("/", pekerjaanService.CreatePekerjaan)
	pekerjaan.Put("/:id", pekerjaanService.UpdatePekerjaan)
	pekerjaan.Delete("/:id", pekerjaanService.DeletePekerjaan)


}