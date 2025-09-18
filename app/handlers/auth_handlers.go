package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"prak4/app/models"
	"prak4/app/utils"
	"prak4/database"
)

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil || req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "username & password wajib"})
	}

	var u models.User
	var hash string
	err := database.DB.QueryRow(`
		SELECT id, username, email, password_hash, role
		FROM users
		WHERE username = $1 OR email = $1
	`, req.Username).Scan(&u.ID, &u.Username, &u.Email, &hash, &u.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(401).JSON(fiber.Map{"error": "username/password salah"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "db error"})
	}

	if !utils.CheckPassword(req.Password, hash) {
		return c.Status(401).JSON(fiber.Map{"error": "username/password salah"})
	}

	tok, err := utils.GenerateToken(u)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal generate token"})
	}

	return c.JSON(models.LoginResponse{User: u, Token: tok})
}

func Profile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"user_id":  c.Locals("user_id"),
		"username": c.Locals("username"),
		"role":     c.Locals("role"),
	})
}
