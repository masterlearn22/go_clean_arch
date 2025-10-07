package service

import (
	"net/mail"
	"strconv" 
    "strings"

	"github.com/gofiber/fiber/v2"
	"go_clean/app/models"
	"go_clean/app/repository"
	"go_clean/utils"
)

type UserService struct {
	Repo *repository.UserRepository
}

func (s *UserService) GetUsersService(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")
	offset := (page - 1) * limit

	// Validasi input
	sortByWhitelist := map[string]bool{"id": true, "name": true,
		"email": true, "created_at": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	// Ambil data dari repository
	users, err := s.Repo.GetUsersRepo(search, sortBy, order,
		limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	total, err := s.Repo.CountUsersRepo(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count users"})
	}

	// Buat response pakai model
	response := models.UserResponse[models.User]{
		Data: users,
		Meta: models.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return c.JSON(response)
}

// helper validasi ringan
func isEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}

// PUBLIC: register user (role = "user" fixed)
func (s *UserService) RegisterUser(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "payload tidak valid"})
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "username, email, password wajib"})
	}
	if !isEmail(req.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "format email tidak valid"})
	}
	exists, err := s.Repo.ExistsByUsernameOrEmail(req.Username, req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "db error"})
	}
	if exists {
		return c.Status(409).JSON(fiber.Map{"error": "username/email sudah dipakai"})
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal hash password"})
	}

	u, err := s.Repo.Create(req.Username, req.Email, hash, "user")
	if err != nil {
		// cek duplikat juga bisa terjadi dari constraint
		return c.Status(500).JSON(fiber.Map{"error": "gagal membuat user"})
	}

	// opsional: langsung login (return token)
	token, err := utils.GenerateToken(*u)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal generate token"})
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "register sukses",
		"user":    u,
		"token":   token,
	})
}

// ADMIN ONLY: create user/admin
func (s *UserService) AdminCreateUser(c *fiber.Ctx) error {
	// pastikan middleware AdminOnly sudah pasang di router
	var req models.AdminCreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "payload tidak valid"})
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.Role = strings.ToLower(strings.TrimSpace(req.Role))

	if req.Username == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		return c.Status(400).JSON(fiber.Map{"error": "username, email, password, role wajib"})
	}
	if !isEmail(req.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "format email tidak valid"})
	}
	if req.Role != "admin" && req.Role != "user" {
		return c.Status(400).JSON(fiber.Map{"error": "role harus 'admin' atau 'user'"})
	}

	exists, err := s.Repo.ExistsByUsernameOrEmail(req.Username, req.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "db error"})
	}
	if exists {
		return c.Status(409).JSON(fiber.Map{"error": "username/email sudah dipakai"})
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal hash password"})
	}

	u, err := s.Repo.Create(req.Username, req.Email, hash, req.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal membuat user"})
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "user dibuat",
		"user":    u,
	})
}
