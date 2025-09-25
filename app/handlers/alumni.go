package handlers

import (
	"github.com/gofiber/fiber/v2"
	"prak4/app/service"
)

func GetAlumniListHandler(c *fiber.Ctx) error {
	return service.GetAlumniList(c)
}
