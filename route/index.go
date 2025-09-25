package route

import (
	"database/sql"

	"prak4/app/handlers"
	"prak4/app/middleware"
	"prak4/app/repository"
	"prak4/app/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	// Initialize repositories
	alumniRepo := &repository.AlumniRepository{DB: db}
	pekerjaanRepo := &repository.PekerjaanRepository{DB: db}

	// Initialize services (tipe method-nya sudah kompatibel: func(*fiber.Ctx) error)
	alumniService := &service.AlumniService{Repo: alumniRepo}
	pekerjaanService := &service.PekerjaanService{Repo: pekerjaanRepo}
	userRepo := &repository.UserRepository{DB: db}
	userService := &service.UserService{Repo: userRepo}
	// Tambahkan inisialisasi userService jika diperlukan dependency lain, tambahkan di sini

	// Root
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Alumni Management API")
	})

	// Base group
	api := app.Group("/api")

	// ========= Public =========
	api.Post("/login", handlers.Login)
	api.Post("/register", userService.RegisterUser)
	api.Get("/alumni-pag", handlers.GetAlumniListHandler)
	api.Get("/pekerjaan-pag",  handlers.GetPekerjaanListHandler)

	// ========= Protected (JWT wajib) =========
	auth := api.Group("", middleware.AuthRequired())
	auth.Group("/users", middleware.AdminOnly())
	auth.Post("/register-admin", middleware.AdminOnly(), userService.AdminCreateUser)

	// Cek profil (ambil dari klaim token)
	auth.Get("/profile", handlers.Profile)

	// ---------- Alumni ----------
	alumni := auth.Group("/alumni")
	// read-only untuk user & admin
	alumni.Get("/", alumniService.GetAllAlumni)
	alumni.Get("/:id", alumniService.GetAlumniByID)
	alumni.Get("/angkatan/:angkatan", alumniService.GetAlumniByAngkatan)

	// write-only untuk admin

	alumniAdmin := alumni.Group("", middleware.AdminOnly())
	alumniAdmin.Post("/", alumniService.CreateAlumni)
	alumniAdmin.Put("/:id", alumniService.UpdateAlumni)
	alumniAdmin.Delete("/:id", alumniService.DeleteAlumni)

	// ---------- Pekerjaan ----------
	pkj := auth.Group("/pekerjaan")
	// read-only untuk user & admin
	pkj.Get("/", pekerjaanService.GetAllPekerjaan)
	pkj.Get("/:id", pekerjaanService.GetPekerjaanByID)
	pkj.Get("/alumni/:alumni_id", pekerjaanService.GetPekerjaanByAlumniID)

	// write-only untuk admin
	pkjAdmin := pkj.Group("", middleware.AdminOnly())
	pkjAdmin.Post("/", pekerjaanService.CreatePekerjaan)
	pkjAdmin.Put("/:id", pekerjaanService.UpdatePekerjaan)
	pkjAdmin.Delete("/:id", pekerjaanService.DeletePekerjaan)
}
