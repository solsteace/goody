package repository

import (
	"time"

	"github.com/solsteace/goody/account/internal/domain"
	"gorm.io/gorm"
)

// Proxy object between persistence layer using Gorm
// and `User` domain object
type gormUser struct {
	ID           uint      `gorm:"column=id"`
	Nama         string    `gorm:"column=nama"`
	KataSandi    string    `gorm:"column=kata_sandi"`
	NoTelp       string    `gorm:"column=no_telp"`
	TanggalLahir time.Time `gorm:"column=tanggal_lahir"`
	JenisKelamin string    `gorm:"column=jenis_kelamin"`
	Tentang      string    `gorm:"column=tentang"`
	Pekerjaan    string    `gorm:"column=pekerjaan"`
	Email        string    `gorm:"column=email"`
	IsAdmin      bool      `gorm:"column=is_admin"`
	IdProvinsi   string    `gorm:"column=id_provinsi"`
	IdKota       string    `gorm:"column=id_kota"`
	UpdatedAt    time.Time `gorm:"column=updated_at"`
	CreatedAt    time.Time `gorm:"column=created_at"`
}

func (gu gormUser) Convert() (domain.User, error) {
	return domain.NewUser(
		gu.Nama,
		gu.KataSandi,
		gu.NoTelp,
		gu.TanggalLahir,
		gu.JenisKelamin,
		gu.Tentang,
		gu.Pekerjaan,
		gu.Email,
		gu.IsAdmin,
		gu.IdProvinsi,
		gu.IdKota,
		gu.UpdatedAt,
		gu.CreatedAt)
}

// @ref https://gorm.io/docs/conventions.html#TableName
func (gu gormUser) TableName() string {
	return "users"
}

func newGormUser(user domain.User) gormUser {
	return gormUser{
		ID:           user.ID,
		Nama:         user.Nama,
		KataSandi:    user.KataSandi,
		NoTelp:       user.NoTelp,
		TanggalLahir: user.TanggalLahir,
		JenisKelamin: user.JenisKelamin,
		Tentang:      user.Tentang,
		Pekerjaan:    user.Pekerjaan,
		Email:        user.Email,
		IsAdmin:      user.IsAdmin,
		IdProvinsi:   user.IdProvinsi,
		IdKota:       user.IdKota,
		UpdatedAt:    user.UpdatedAt,
		CreatedAt:    user.CreatedAt,
	}
}

type gormUserRepo struct {
	db *gorm.DB
}

func NewGormUserRepo(db *gorm.DB) gormUserRepo {
	return gormUserRepo{db: db}
}

func (gur gormUserRepo) Migrate() {
	gur.db.AutoMigrate(new(gormUser))
}

func (gur gormUserRepo) GetById(id int) (Userable, error) {
	user := new(gormUser)
	result := gur.db.
		Where("id = ?", id).
		First(&user)

	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (gur gormUserRepo) GetByPhoneNumber(phone string) (Userable, error) {
	user := new(gormUser)
	result := gur.db.
		Where("no_telp = ?", phone).
		First(&user)

	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (gur gormUserRepo) Create(u domain.User) (uint, error) {
	user := newGormUser(u)
	result := gur.db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}
