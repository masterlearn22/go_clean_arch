package service

import (
	"database/sql"
	"fmt"
	"go_clean/app/models"
	"go_clean/app/repository"
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
			Pages:  (total + params.Limit - 1) / params.Limit,
			SortBy: params.SortBy, Order: params.Order, Search: params.Search,
		},
	}
	return c.JSON(resp)
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

func (s *AlumniService) GetAlumniAndPekerjaan(c *fiber.Ctx) error {
	idStr := c.Params("nim") // sebenarnya ini ID
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID harus berupa angka",
		})
	}

	result, err := s.Repo.GetAlumniAndPekerjaan(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data alumni dan pekerjaan: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data alumni dan pekerjaan berhasil diambil",
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
