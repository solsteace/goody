package domain

import (
	"time"
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
	// TODO: domain validation...
	user := User{
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

func (u User) WithId(id uint) User {
	u.ID = id
	return u
}
