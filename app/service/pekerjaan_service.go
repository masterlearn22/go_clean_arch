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
    // --- Ambil parameter ID pekerjaan ---
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "ID pekerjaan tidak valid",
        })
    }

    // --- Ambil data user dari token ---
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(int)

    // --- Ambil user dari DB untuk tahu alumni_id ---
    userRepo := repository.UserRepository{DB: s.Repo.DB}
    user, err := userRepo.GetUserByID(userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Gagal mengambil data user: " + err.Error(),
        })
    }

    // --- Ambil data pekerjaan lama ---
    existing, err := s.Repo.GetPekerjaanByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error":   true,
            "message": "Data pekerjaan tidak ditemukan",
        })
    }

    // --- Validasi kepemilikan data ---
    if role == "user" && user.AlumniID != nil && existing.AlumniID != *user.AlumniID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error":   true,
            "message": "Kamu tidak punya izin mengubah pekerjaan ini",
        })
    }

    // --- Parse request body ---
    var p models.PekerjaanAlumni
    if err := c.BodyParser(&p); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   true,
            "message": "Request body tidak valid",
        })
    }

    // --- Set alumni_id otomatis untuk user biasa ---
    if role == "user" {
        p.AlumniID = *user.AlumniID
    }

    // --- Update ke database ---
    rows, err := s.Repo.UpdatePekerjaan(id, &p)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   true,
            "message": "Gagal mengupdate pekerjaan: " + err.Error(),
        })
    }

    if rows == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error":   true,
            "message": "Pekerjaan tidak ditemukan untuk diupdate",
        })
    }

    updated, _ := s.Repo.GetPekerjaanByID(id)
    return c.JSON(fiber.Map{
        "success": true,
        "message": "Pekerjaan berhasil diupdate",
        "data":    updated,
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

    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(int)

    // Ambil alumni_id dari user login
    userRepo := repository.UserRepository{DB: s.Repo.DB}
    user, err := userRepo.GetUserByID(userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal mengambil data user: " + err.Error(),
        })
    }

    // Ambil data pekerjaan
    existing, err := s.Repo.GetPekerjaanByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Pekerjaan tidak ditemukan",
        })
    }

    // Validasi kepemilikan
    if role == "user" && user.AlumniID != nil && existing.AlumniID != *user.AlumniID  {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "success": false,
            "message": "Kamu tidak memiliki izin menghapus pekerjaan ini",
        })
    }

    // Soft delete (gunakan deleted_by sesuai user login)
    rows, err := s.Repo.SoftDeletePekerjaan(id, userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal menghapus pekerjaan: " + err.Error(),
        })
    }

    if rows == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Pekerjaan tidak ditemukan untuk dihapus",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Pekerjaan berhasil dihapus (soft delete)",
    })
}



func (s *PekerjaanService) TrashAllPekerjaan(c *fiber.Ctx) error {
    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(int)

    // Ambil alumni_id dari tabel users
    userRepo := repository.UserRepository{DB: s.Repo.DB}
    user, err := userRepo.GetUserByID(userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal mengambil data user: " + err.Error(),
        })
    }

    var pekerjaan []models.PekerjaanAlumni

    // 🔑 Role-based access
    if role == "admin" {
        // Admin bisa lihat semua data di trash
        pekerjaan, err = s.Repo.TrashAllPekerjaan()
    } else {
        // User hanya bisa lihat data miliknya
        pekerjaan, err = s.Repo.TrashPekerjaanByAlumniID(*user.AlumniID)
    }

    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal mengambil data pekerjaan: " + err.Error(),
        })
    }

    if len(pekerjaan) == 0 {
        return c.JSON(fiber.Map{
            "success": true,
            "message": "Tidak ada data pekerjaan di trash",
            "data":    []interface{}{},
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
    userID := c.Locals("user_id").(int)
    pekerjaanID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "ID pekerjaan tidak valid",
        })
    }

    userRepo := repository.UserRepository{DB: s.Repo.DB}
    user, err := userRepo.GetUserByID(userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal mengambil data user: " + err.Error(),
        })
    }

    existing, err := s.Repo.GetPekerjaanByID(pekerjaanID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Data pekerjaan tidak ditemukan",
        })
    }

    if role == "user" && user.AlumniID != nil && existing.AlumniID != *user.AlumniID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "success": false,
            "message": "Kamu tidak punya izin untuk me-restore pekerjaan ini",
        })
    }

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
    pekerjaanID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "ID pekerjaan tidak valid",
        })
    }

    userRepo := repository.UserRepository{DB: s.Repo.DB}
    user, err := userRepo.GetUserByID(userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Gagal mengambil data user: " + err.Error(),
        })
    }

    existing, err := s.Repo.GetPekerjaanByID(pekerjaanID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Pekerjaan tidak ditemukan",
        })
    }

    if role == "user" && user.AlumniID != nil && existing.AlumniID != *user.AlumniID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "success": false,
            "message": "Kamu tidak memiliki izin menghapus permanen pekerjaan ini",
        })
    }

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