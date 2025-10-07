package repository

import (
	"fmt"
	"database/sql"
	"go_clean/app/models"
	"time"
	"go_clean/database"
)

type AlumniRepository struct {
	DB *sql.DB
}
// AlumniSortable returns a slice of sortable field names for alumni
func AlumniSortable() []string {
    return []string{"nim", "nama", "jurusan", "angkatan", "email"}
}

func (r *AlumniRepository) GetAllAlumni() ([]models.Alumni, error) {
	rows, err := r.DB.Query("SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []models.Alumni
	for rows.Next() {
		var a models.Alumni
		if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		alumniList = append(alumniList, a)
	}
	return alumniList, nil
}

func (r *AlumniRepository) GetAlumniByID(id int) (*models.Alumni, error) {
	var a models.Alumni
	err := r.DB.QueryRow("SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at FROM alumni WHERE id = $1", id).Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AlumniRepository) GetAlumniByAngkatan(angkatan int) (*models.AlumniAngkatan, error) {
    jumlahalumni := &models.AlumniAngkatan{Angkatan: angkatan}
    err := r.DB.QueryRow("SELECT COUNT(*) FROM alumni WHERE angkatan = $1", angkatan).Scan(&jumlahalumni.Jumlah)
    if err != nil {
        return nil, err
    }
    return jumlahalumni, nil
}

func (r *AlumniRepository) CreateAlumni(alumni *models.Alumni) (int, error) {
	var id int
	err := r.DB.QueryRow(
		"INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		alumni.NIM, alumni.Nama, alumni.Jurusan, alumni.Angkatan, alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat, time.Now(), time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AlumniRepository) UpdateAlumni(id int, alumni *models.Alumni) (int64, error) {
	result, err := r.DB.Exec(
		"UPDATE alumni SET nama = $1, jurusan = $2, angkatan = $3, tahun_lulus = $4, email = $5, no_telepon = $6, alamat = $7, updated_at = $8 WHERE id = $9",
		alumni.Nama, alumni.Jurusan, alumni.Angkatan, alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat, time.Now(), id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *AlumniRepository) DeleteAlumni(id int) (int64, error) {
	result, err := r.DB.Exec("DELETE FROM alumni WHERE id = $1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func sanitizeAlumniSort(s string) string {
    switch s {
    case "id", "nim", "nama", "jurusan", "angkatan", "email", "created_at", "updated_at":
        return s
    default:
        return "id"
    }
}

func sanitizeOrderAlumni(o string) string {
    if o == "desc" || o == "DESC" {
        return "DESC"
    }
    return "ASC"
}


func ListAlumniRepo(search, sortBy, order string, limit, offset int) ([]models.Alumni, error) {
    // Sanitasi sort & order biar aman dari SQL injection via fmt.Sprintf
    sortBy = sanitizeAlumniSort(sortBy)
    order  = sanitizeOrderAlumni(order)

    query := fmt.Sprintf(`
        SELECT id, nama, nim, angkatan
        FROM alumni
        WHERE (nama ILIKE $1 OR CAST(nim AS TEXT) ILIKE $1)
        ORDER BY %s %s, id ASC
        LIMIT $2 OFFSET $3
    `, sortBy, order)

    rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.Alumni
    for rows.Next() {
        var a models.Alumni
        if err := rows.Scan(&a.ID, &a.Nama, &a.NIM, &a.Angkatan); err != nil {
            return nil, err
        }
        items = append(items, a)
    }
    return items, rows.Err()
}

func CountAlumniRepo(search string) (int, error) {
    var total int
    err := database.DB.QueryRow(`
        SELECT COUNT(*)
        FROM alumni
        WHERE (nama ILIKE $1 OR CAST(nim AS TEXT) ILIKE $1)
    `, "%"+search+"%").Scan(&total)
    if err != nil && err != sql.ErrNoRows {
        return 0, err
    }
    return total, nil
}