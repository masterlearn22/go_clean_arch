package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlumniMongo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AlumniID    int                `bson:"alumni_id" json:"alumni_id"`
	NIM         string             `bson:"nim" json:"nim"`
	Nama        string             `bson:"nama" json:"nama"`
	Jurusan     string             `bson:"jurusan" json:"jurusan"`
	Angkatan    int                `bson:"angkatan" json:"angkatan"`
	TahunLulus  int                `bson:"tahun_lulus" json:"tahun_lulus"`
	Email       string             `bson:"email" json:"email"`
	NoTelp      string             `bson:"no_telepon" json:"no_telepon"`
	Alamat      string             `bson:"alamat" json:"alamat"`
	TempatKerja string             `bson:"tempat_kerja,omitempty" json:"tempat_kerja,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
