package service

import (
	"fmt"
	"fmt"
	"database/sql"
	"prak4/app/models"
	"prak4/app/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	Repo *repository.AlumniRepository
}

func (s *AlumniService) GetAllAlumni(c *fiber.Ctx) error {
	alumni, err := s.Repo.GetAllAlumni()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data alumni: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil",
		"data":    alumni,
	})
}

func GetAlumniList(c *fiber.Ctx) error {
	sortable := make(map[string]bool)
	for _, v := range repository.AlumniSortable() {
		sortable[v] = true
	}
params := getListParams(c, sortable) // lihat fungsi accessor di bawah
items, err := repository.ListAlumniRepo(params.Search, params.SortBy, params.Order, params.Limit, params.Offset)
if err != nil {
	// DEBUG: log detail error
	fmt.Printf("ListAlumniRepo error: %v\n", err)
	return c.Status(500).JSON(fiber.Map{
		"error": "failed to fetch alumni",
	})
}

	total, err := repository.CountAlumniRepo(params.Search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to count alumni"})
	}

	resp := models.UserResponse[models.Alumni]{
		Data: items,
		Meta: models.MetaInfo{
			Page: params.Page, Limit: params.Limit, Total: total,
			Pages: (total + params.Limit - 1) / params.Limit,
			SortBy: params.SortBy, Order: params.Order, Search: params.Search,
		},
	}
	return c.JSON(resp)
}

// getListParams extracts pagination, sorting, and search parameters from the request context.
func getListParams(c *fiber.Ctx, sortable map[string]bool) struct {
	Page   int
	Limit  int
	Offset int
	SortBy string
	Order  string
	Search string
} {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	sortBy := c.Query("sortBy", "id")
	if !sortable[sortBy] {
		sortBy = "id"
	}
	order := c.Query("order", "asc")
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	search := c.Query("search", "")

	return struct {
		Page   int
		Limit  int
		Offset int
		SortBy string
		Order  string
		Search string
	}{
		Page:   page,
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
		Order:  order,
		Search: search,
	}
}

func (s *AlumniService) GetAlumniByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	alumni, err := s.Repo.GetAlumniByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Alumni tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data alumni: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil",
		"data":    alumni,
	})
}

func (s *AlumniService) GetAlumniByAngkatan(c *fiber.Ctx) error {
    angkatan, err := strconv.Atoi(c.Params("angkatan"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Angkatan tidak valid",
        })
    }

    result, err := s.Repo.GetAlumniByAngkatan(angkatan)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal mengambil data alumni: " + err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Data alumni berhasil diambil",
        "data":    result,
    })
}

func (s *AlumniService) CreateAlumni(c *fiber.Ctx) error {
	var alumni models.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	if alumni.NIM == "" || alumni.Nama == "" || alumni.Jurusan == "" || alumni.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field NIM, Nama, Jurusan, dan Email wajib diisi",
		})
	}

	newID, err := s.Repo.CreateAlumni(&alumni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menambah alumni: " + err.Error(),
		})
	}

	newAlumni, _ := s.Repo.GetAlumniByID(newID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil ditambahkan",
		"data":    newAlumni,
	})
}

func (s *AlumniService) UpdateAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	var alumni models.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	rowsAffected, err := s.Repo.UpdateAlumni(id, &alumni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengupdate alumni: " + err.Error(),
		})
	}
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Alumni tidak ditemukan untuk diupdate",
		})
	}

	updatedAlumni, _ := s.Repo.GetAlumniByID(id)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil diupdate",
		"data":    updatedAlumni,
	})
}

func (s *AlumniService) DeleteAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	rowsAffected, err := s.Repo.DeleteAlumni(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghapus alumni: " + err.Error(),
		})
	}
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Alumni tidak ditemukan untuk dihapus",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil dihapus",
	})
}