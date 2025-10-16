package models

import "time"

// PekerjaanAlumni merepresentasikan tabel pekerjaan_alumni
type PekerjaanAlumni struct {
	ID                 int        `json:"id"`
	AlumniID           int        `json:"alumni_id"`
	NamaPerusahaan     string     `json:"nama_perusahaan"`
	PosisiJabatan      string     `json:"posisi_jabatan"`
	BidangIndustri     string     `json:"bidang_industri"`
	LokasiKerja        string     `json:"lokasi_kerja"`
	GajiRange          *string    `json:"gaji_range"`
	TanggalMulaiKerja  time.Time  `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan    string     `json:"status_pekerjaan"`
	DeskripsiPekerjaan *string    `json:"deskripsi_pekerjaan"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	IsDeleted          bool       `json:"is_deleted"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	DeletedBy string      `json:"deleted_by"`
}


type TrashPekerjaanAlumni struct {
	ID                 int        `json:"id"`
	AlumniID           int        `json:"alumni_id"`
	NamaPerusahaan     string     `json:"nama_perusahaan"`
	PosisiJabatan      string     `json:"posisi_jabatan"`
	BidangIndustri     string     `json:"bidang_industri"`
	LokasiKerja        string     `json:"lokasi_kerja"`
	GajiRange          *string    `json:"gaji_range"`
	TanggalMulaiKerja  time.Time  `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan    string     `json:"status_pekerjaan"`
	DeskripsiPekerjaan *string    `json:"deskripsi_pekerjaan"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	IsDeleted          bool       `json:"is_deleted"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	DeletedBy string      `json:"deleted_by"`
}