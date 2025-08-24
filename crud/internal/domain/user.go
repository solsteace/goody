package domain

import (
	"errors"
	"time"

	"github.com/solsteace/goody/lib/oops"
)

type User struct {
	ID           uint      `json:"id"`
	Nama         string    `json:"nama"`
	KataSandi    string    `json:"kata_sandi"`
	NoTelp       string    `json:"no_telp"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	JenisKelamin string    `json:"jenis_kelamin"`
	Tentang      string    `json:"tentang"`
	Pekerjaan    string    `json:"pekerjaan"`
	Email        string    `json:"email"`
	IsAdmin      bool      `json:"is_admin"`
	IdProvinsi   string    `json:"id_provinsi"`
	IdKota       string    `json:"id_kota"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewUser(
	id *uint,
	nama,
	kataSandi,
	noTelp string,
	tanggalLahir time.Time,
	jenisKelamin,
	tentang,
	pekerjaan,
	email string,
	isAdmin bool,
	idProvinsi,
	idKota string,
	updatedAt time.Time,
	createdAt time.Time,
) (User, error) {
	var userId uint = 0
	if id != nil {
		userId = *id
	}

	switch jenisKelamin {
	case "pria", "wanita", "anonim":
	default:
		return User{}, oops.BadValues{
			Err: errors.New("`JenisKelamin` value invalid: " + jenisKelamin),
			Msg: "Jenis kelamin harus salah satu di antara 'pria', 'wanita', atau 'anonim'"}
	}

	user := User{
		ID:           userId,
		Nama:         nama,
		KataSandi:    kataSandi,
		NoTelp:       noTelp,
		TanggalLahir: tanggalLahir,
		JenisKelamin: jenisKelamin,
		Tentang:      tentang,
		Pekerjaan:    pekerjaan,
		Email:        email,
		IsAdmin:      isAdmin,
		IdProvinsi:   idProvinsi,
		IdKota:       idKota,
		UpdatedAt:    updatedAt,
		CreatedAt:    createdAt,
	}
	return user, nil
}
