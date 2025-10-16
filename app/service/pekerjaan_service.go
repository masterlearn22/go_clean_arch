package service

import (
	"fmt"
	"database/sql"
	"go_clean/app/models"
	"go_clean/app/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanService struct {
	Repo *repository.PekerjaanRepository
}

// Ambil semua pekerjaan tanpa filter/pagination
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

// Ambil pekerjaan berdasarkan ID
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

// Ambil list pekerjaan dengan search, sort, pagination (mirip AlumniService)
func GetPekerjaanList(c *fiber.Ctx) error {
	sortable := repository.PekerjaanSortable()
	params := getListParams(c, sortable)

	items, err := repository.ListPekerjaanRepo(params.Search, params.SortBy, params.Order, params.Limit, params.Offset)
	if err != nil {
		fmt.Printf("ListPekerjaanRepo error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch pekerjaan"})
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

// Ambil semua pekerjaan milik alumni tertentu
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

// Tambah data pekerjaan baru
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

// Update data pekerjaan berdasarkan ID
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
			"message": "Pekerjaan tidak ditemukan untuk diupdate",
		})
	}

	updatedPekerjaan, _ := s.Repo.GetPekerjaanByID(id)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil diupdate",
		"data":    updatedPekerjaan,
	})
}

// Hapus pekerjaan (soft delete)
func (s *PekerjaanService) DeletePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID pekerjaan tidak valid",
		})
	}

	// Kalau DeletePekerjaan di repository butuh deletedBy, tambahkan parameter user ID di sini
	rowsAffected, err := s.Repo.SoftDeletePekerjaan(id, 0)
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

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil dihapus",
	})
}


func (s *PekerjaanService) TrashAllPekerjaan(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	var pekerjaan interface{}
	var err error

	if role == "admin" {
		pekerjaan, err = s.Repo.TrashAllPekerjaan()
	} else {
		pekerjaan, err = s.Repo.TrashPekerjaanByAlumniID(userID)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan trash berhasil diambil",
		"data":    pekerjaan,
	})
}

func (s *PekerjaanService) RestorePekerjaan(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int) // ini alumni_id dari token
	idStr := c.Params("id")

	pekerjaanID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID pekerjaan tidak valid",
		})
	}

	// Kalau bukan admin, pastikan pekerjaan ini milik user
	if role != "admin" {
		isOwner, err := s.Repo.IsPekerjaanOwnedByUser(pekerjaanID, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal memeriksa kepemilikan pekerjaan: " + err.Error(),
			})
		}
		if !isOwner {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Anda tidak memiliki pekerjaan ini",
			})
		}
	}

	// Jalankan restore
	err = s.Repo.RestorePekerjaanByID(pekerjaanID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal me-restore pekerjaan: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil di-restore",
	})
}


func (s *PekerjaanService) HardDeletePekerjaan(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)
	idStr := c.Params("id")

	pekerjaanID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID pekerjaan tidak valid",
		})
	}

	// Kalau user biasa, pastikan pekerjaan itu miliknya dan sudah di trash
	if role != "admin" {
		isOwnerAndTrashed, err := s.Repo.IsTrashedPekerjaanOwnedByUser(pekerjaanID, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Gagal memeriksa kepemilikan pekerjaan: " + err.Error(),
			})
		}
		if !isOwnerAndTrashed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Anda tidak memiliki pekerjaan ini atau pekerjaan belum di-trash",
			})
		}
	}

	// Lanjut hard delete
	err = s.Repo.HardDeletePekerjaanByID(pekerjaanID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghapus permanen pekerjaan: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil dihapus permanen",
	})
}

