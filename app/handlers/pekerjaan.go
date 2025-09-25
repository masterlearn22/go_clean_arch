package handlers

import (
	"github.com/gofiber/fiber/v2"
	"prak4/app/service"
)

func GetPekerjaanListHandler(c *fiber.Ctx) error {
	return service.GetPekerjaanList(c)
}
