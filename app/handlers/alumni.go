package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go_clean/app/service"
)

func GetAlumniListHandler(c *fiber.Ctx) error {
	return service.GetAlumniList(c)
}
