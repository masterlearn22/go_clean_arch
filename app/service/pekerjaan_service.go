package service

import (
	"fmt"
	"database/sql"
	"fmt"
	"prak4/app/models"
	"prak4/app/repository"
	"prak4/helper"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanService struct {
	Repo *repository.PekerjaanRepository
}

func (s *PekerjaanService) GetAllPekerjaan(c *fiber.Ctx) error {
	pekerjaan, err := s.Repo.GetAllPekerjaan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan berhasil diambil",
		"data":    pekerjaan,
	})
}

func (s *PekerjaanService) GetPekerjaanByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID pekerjaan tidak valid",
		})
	}

	pekerjaan, err := s.Repo.GetPekerjaanByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pekerjaan tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan berhasil diambil",
		"data":    pekerjaan,
	})
}

func (s *PekerjaanService) GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID alumni tidak valid",
		})
	}

	pekerjaan, err := s.Repo.GetPekerjaanByAlumniID(alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan untuk alumni berhasil diambil",
		"data":    pekerjaan,
	})
}

func GetPekerjaanList(c *fiber.Ctx) error {
	sortable := repository.PekerjaanSortable()
	params := getListParams(c, sortable)

	items, err := repository.ListPekerjaanRepo(params.Search, params.SortBy, params.Order, params.Limit, params.Offset)
	if err != nil {
		// DEBUG: log detail error
		fmt.Printf("ListPekerjaanRepo error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch pekerjaan",
		})
	}

	total, err := repository.CountPekerjaanRepo(params.Search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to count pekerjaan"})
	}

	resp := models.UserResponse[models.PekerjaanAlumni]{
		Data: items,
		Meta: models.MetaInfo{
			Page: params.Page, Limit: params.Limit, Total: total,
			Pages:  (total + params.Limit - 1) / params.Limit,
			SortBy: params.SortBy, Order: params.Order, Search: params.Search,
		},
	}
	return c.JSON(resp)
}

func (s *PekerjaanService) CreatePekerjaan(c *fiber.Ctx) error {
	var p models.PekerjaanAlumni
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	if p.AlumniID == 0 || p.NamaPerusahaan == "" || p.PosisiJabatan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field alumni_id, nama_perusahaan, dan posisi_jabatan wajib diisi",
		})
	}

	newID, err := s.Repo.CreatePekerjaan(&p)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menambah pekerjaan: " + err.Error(),
		})
	}

	newPekerjaan, _ := s.Repo.GetPekerjaanByID(newID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil ditambahkan",
		"data":    newPekerjaan,
	})
}

func (s *PekerjaanService) UpdatePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	var p models.PekerjaanAlumni
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	rowsAffected, err := s.Repo.UpdatePekerjaan(id, &p)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengupdate pekerjaan: " + err.Error(),
		})
	}
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Pekerjaan tidak ditemukan",
		})
	}

	updatedPekerjaan, _ := s.Repo.GetPekerjaanByID(id)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil diupdate",
		"data":    updatedPekerjaan,
	})
}

func (s *PekerjaanService) DeletePekerjaan(c *fiber.Ctx) error {
    // Ambil ID pekerjaan dari parameter URL
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusBadRequest, "ID pekerjaan tidak valid")
    }

<<<<<<< HEAD
	rowsAffected, err := s.Repo.DeletePekerjaan(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghapus pekerjaan: " + err.Error(),
		})
	}
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Pekerjaan tidak ditemukan",
		})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus"})
=======
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(int)

    pekerjaan, err := s.Repo.GetPekerjaanByID(id)
    if err != nil {
        if err == sql.ErrNoRows {
            return helper.ErrorResponse(c, fiber.StatusNotFound, "Pekerjaan tidak ditemukan")
        }
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data pekerjaan")
    }

    // ❌ Bukan admin & bukan pemilik pekerjaan
    if role != "admin" && pekerjaan.AlumniID != userID {
        return helper.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak: Anda hanya dapat menghapus pekerjaan Anda sendiri.")
    }

    // ✅ Soft delete: set is_delete + deleted_at + deleted_by
    rowsAffected, err := s.Repo.DeletePekerjaan(id, userID)
    if err != nil {
        return helper.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal menghapus pekerjaan: "+err.Error())
    }
    if rowsAffected == 0 {
        return helper.ErrorResponse(c, fiber.StatusNotFound, "Pekerjaan tidak ditemukan untuk dihapus")
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": fmt.Sprintf("Pekerjaan berhasil dihapus oleh user_id=%d", userID),
    })
>>>>>>> af843f7 (Memperbarui struktur porject)
}
