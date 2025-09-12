package service

import (
	"database/sql"
	"prak4/app/model"
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
	var alumni model.Alumni
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

	var alumni model.Alumni
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