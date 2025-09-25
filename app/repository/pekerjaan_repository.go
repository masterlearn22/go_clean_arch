package repository

import (
	"fmt"
	"database/sql"
	"prak4/app/models"
	"time"
	"prak4/database"
)

type PekerjaanRepository struct {
	DB *sql.DB
}

func PekerjaanSortable() map[string]bool {
	return map[string]bool{
		"id": true, "alumni_id": true, "nama_perusahaan": true, "posisi_jabatan": true, "tanggal_mulai_kerja": true,
	}
}

func (r *PekerjaanRepository) GetAllPekerjaan() ([]models.PekerjaanAlumni, error) {
	rows, err := r.DB.Query("SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at FROM pekerjaan_alumni ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []models.PekerjaanAlumni
	for rows.Next() {
		var p models.PekerjaanAlumni
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func (r *PekerjaanRepository) GetPekerjaanByID(id int) (*models.PekerjaanAlumni, error) {
	var p models.PekerjaanAlumni
	err := r.DB.QueryRow("SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at FROM pekerjaan_alumni WHERE id = $1", id).Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PekerjaanRepository) GetPekerjaanByAlumniID(alumniID int) ([]models.PekerjaanAlumni, error) {
	rows, err := r.DB.Query("SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at FROM pekerjaan_alumni WHERE alumni_id = $1 ORDER BY tanggal_mulai_kerja DESC", alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []models.PekerjaanAlumni
	for rows.Next() {
		var p models.PekerjaanAlumni
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func (r *PekerjaanRepository) CreatePekerjaan(p *models.PekerjaanAlumni) (int, error) {
	var id int
	err := r.DB.QueryRow(
		`INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
		p.AlumniID, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange, p.TanggalMulaiKerja, p.TanggalSelesaiKerja, p.StatusPekerjaan, p.DeskripsiPekerjaan, time.Now(), time.Now(),
	).Scan(&id)
	return id, err
}

func (r *PekerjaanRepository) UpdatePekerjaan(id int, p *models.PekerjaanAlumni) (int64, error) {
	result, err := r.DB.Exec(
		`UPDATE pekerjaan_alumni SET nama_perusahaan = $1, posisi_jabatan = $2, bidang_industri = $3, lokasi_kerja = $4, gaji_range = $5, tanggal_mulai_kerja = $6, tanggal_selesai_kerja = $7, status_pekerjaan = $8, deskripsi_pekerjaan = $9, updated_at = $10 
		 WHERE id = $11`,
		p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange, p.TanggalMulaiKerja, p.TanggalSelesaiKerja, p.StatusPekerjaan, p.DeskripsiPekerjaan, time.Now(), id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *PekerjaanRepository) DeletePekerjaan(id int) (int64, error) {
	result, err := r.DB.Exec("DELETE FROM pekerjaan_alumni WHERE id = $1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}


func ListPekerjaanRepo(search, sortBy, order string, limit, offset int) ([]models.PekerjaanAlumni, error) {
	query := fmt.Sprintf(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, tanggal_mulai_kerja
		FROM pekerjaan_alumni
		WHERE (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1)
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.PekerjaanAlumni
	for rows.Next() {
		var p models.PekerjaanAlumni
		if err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.TanggalMulaiKerja); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, rows.Err()
}

func CountPekerjaanRepo(search string) (int, error) {
	var total int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM pekerjaan_alumni
		WHERE (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1)
	`, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}