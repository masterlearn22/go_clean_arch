package route

import (
	"database/sql"

	"go_clean/app/handlers"
	"go_clean/app/repository"
	"go_clean/app/service"
	"go_clean/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	// =======================
	// REPOSITORIES
	// =======================
	alumniRepo := &repository.AlumniRepository{DB: db}
	pekerjaanRepo := &repository.PekerjaanRepository{DB: db}
	userRepo := &repository.UserRepository{DB: db}

	// =======================
	// SERVICES
	// =======================
	alumniService := &service.AlumniService{Repo: alumniRepo}
	pekerjaanService := &service.PekerjaanService{Repo: pekerjaanRepo}
	userService := &service.UserService{Repo: userRepo}

	// =======================
	// ROOT
	// =======================
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Alumni Management API 🚀")
	})

	// =======================
	// PUBLIC
	// =======================
	api := app.Group("/api")
	api.Post("/login", handlers.Login)
	api.Post("/register", userService.RegisterUser)

	// =======================
	// PROTECTED (JWT REQUIRED)
	// =======================
	auth := api.Group("", middleware.AuthRequired())

	// -- Admin create user --
	auth.Post("/register-admin", middleware.AdminOnly(), userService.AdminCreateUser)

	// -- Profile user login --
	auth.Get("/profile", handlers.Profile)

	// =======================
	// ALUMNI ROUTES
	// =======================
	alumni := auth.Group("/alumni")

	// ✅ Read-only untuk semua user login
	alumni.Get("/", alumniService.GetAllAlumni)
	alumni.Get("/:id", alumniService.GetAlumniByID)
	alumni.Get("/angkatan/:angkatan", alumniService.GetAlumniByAngkatan)
	alumni.Get("/alumni-pag", handlers.GetAlumniListHandler)
	alumni.Get("/with-pekerjaan/:nim", alumniService.GetAlumniAndPekerjaan)

	// ✅ Write-only khusus admin
	alumniAdmin := alumni.Group("", middleware.AdminOnly())
	alumniAdmin.Post("/", alumniService.CreateAlumni)
	alumniAdmin.Put("/:id", alumniService.UpdateAlumni)
	alumniAdmin.Delete("/:id", alumniService.DeleteAlumni)

	// =======================
	// PEKERJAAN ROUTES
	// =======================
	pkj := auth.Group("/pekerjaan")

	// ✅ Semua user login bisa baca
	pkj.Get("/trash", pekerjaanService.TrashAllPekerjaan)
	pkj.Get("/", pekerjaanService.GetAllPekerjaan)
	pkj.Get("/:id", pekerjaanService.GetPekerjaanByID)
	pkj.Get("/alumni/:alumni_id", pekerjaanService.GetPekerjaanByAlumniID)
	pkj.Put("/restore/:id", pekerjaanService.RestorePekerjaan)
	pkj.Delete("/hard-delete/:id", pekerjaanService.HardDeletePekerjaan)

	// ✅ Semua user login bisa update & delete
	//    (validasi “milik sendiri” dilakukan di service)
	pkj.Put("/:id", pekerjaanService.UpdatePekerjaan)
	pkj.Delete("/:id", pekerjaanService.DeletePekerjaan)

	// ✅ Admin-only untuk Create
	pkjAdmin := pkj.Group("", middleware.AdminOnly())
	pkjAdmin.Post("/", pekerjaanService.CreatePekerjaan)

	// =======================
	// PAGINATION HANDLER
	// =======================
	api.Get("/pekerjaan-pag", handlers.GetPekerjaanListHandler)
}
