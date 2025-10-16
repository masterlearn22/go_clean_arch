package models

import "time"

// Alumni merepresentasikan tabel alumni di database
type Alumni struct {
	ID         int       `json:"id"`
	NIM        string    `json:"nim"`
	Nama       string    `json:"nama"`
	Jurusan    string    `json:"jurusan"`
	Angkatan   int       `json:"angkatan"`
	TahunLulus int       `json:"tahun_lulus"`
	Email      string    `json:"email"`
	NoTelepon  *string   `json:"no_telepon"`
	Alamat     *string   `json:"alamat"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AlumniAngkatan struct {
	Angkatan int `json:"angkatan"`
	Jumlah   int `json:"jumlah"`
}

type AlumniPekerjaan struct {
    ID             int    `json:"id"`
    NIM            string `json:"nim"`
    Nama           string `json:"nama"`
    Jurusan        string `json:"jurusan"`
    Angkatan       int    `json:"angkatan"`
    TahunLulus     int    `json:"tahun_lulus"`
    Email          string `json:"email"`
    NamaPerusahaan string `json:"nama_perusahaan"`
    Posisi         string `json:"posisi_jabatan"`
    TahunMulai     int    `json:"tanggal_mulai_kerja"`
    TahunSelesai   int    `json:"tanggal_selesai_kerja"`
}
