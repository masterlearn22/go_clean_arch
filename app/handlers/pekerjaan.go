package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go_clean/app/service"
)

func GetPekerjaanListHandler(c *fiber.Ctx) error {
	return service.GetPekerjaanList(c)
}
