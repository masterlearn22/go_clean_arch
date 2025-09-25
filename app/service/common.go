package service

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ListParams struct {
	Page   int
	Limit  int
	SortBy string
	Order  string
	Search string
	Offset int
}

func getListParams(c *fiber.Ctx, whitelist map[string]bool) ListParams {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 { page = 1 }

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 { limit = 10 }
	if limit > 100 { limit = 100 }

	sortBy := c.Query("sortBy", "id")
	if !whitelist[sortBy] {
		sortBy = "id"
	}

	order := strings.ToLower(c.Query("order", "asc"))
	if order != "desc" {
		order = "asc"
	}

	search := c.Query("search", "")

	return ListParams{
		Page: page, Limit: limit, SortBy: sortBy, Order: order, Search: search,
		Offset: (page - 1) * limit,
	}
}
