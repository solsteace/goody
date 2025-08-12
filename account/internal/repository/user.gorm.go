package repository

import (
	"time"

	"github.com/solsteace/goody/account/internal/domain"
	"github.com/solsteace/goody/account/internal/lib/errors"
	"gorm.io/gorm"
)

// Proxy object between persistence layer using Gorm
// and `User` domain object
type gormUser struct {
	ID           uint      `gorm:"column:id"`
	Nama         string    `gorm:"column:nama"`
	KataSandi    string    `gorm:"column:kata_sandi"`
	NoTelp       string    `gorm:"column:no_telp"`
	TanggalLahir time.Time `gorm:"column:tanggal_lahir"`
	JenisKelamin string    `gorm:"column:jenis_kelamin"`
	Tentang      string    `gorm:"column:tentang"`
	Pekerjaan    string    `gorm:"column:pekerjaan"`
	Email        string    `gorm:"column:email"`
	IsAdmin      bool      `gorm:"column:is_admin"`
	IdProvinsi   string    `gorm:"column:id_provinsi"`
	IdKota       string    `gorm:"column:id_kota"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

// @ref https://gorm.io/docs/conventions.html#TableName
func (gu gormUser) TableName() string {
	return "users"
}

func (gu gormUser) toUser() (domain.User, error) {
	user, err := domain.NewUser(
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
	if err != nil {
		return domain.User{}, err
	}

	return user.WithId(gu.ID), nil
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

func (gur gormUserRepo) GetById(id uint) (domain.User, error) {
	row := new(gormUser)
	result := gur.db.
		Where("id = ?", id).
		First(&row)
	if result.Error != nil {
		return domain.User{}, errors.Standardize(result.Error)
	}

	user, err := row.toUser()
	if err != nil {
		return domain.User{}, errors.Standardize(err)
	}
	return user, errors.Standardize(err)
}

func (gur gormUserRepo) GetByPhoneNumber(phone string) (domain.User, error) {
	row := new(gormUser)
	result := gur.db.
		Where("no_telp = ?", phone).
		First(&row)
	if result.Error != nil {
		return domain.User{}, errors.Standardize(result.Error)
	}

	user, err := row.toUser()
	if err != nil {
		return domain.User{}, errors.Standardize(result.Error)
	}
	return user, errors.Standardize(result.Error)
}

func (gur gormUserRepo) Create(u domain.User) (uint, error) {
	user := newGormUser(u)
	result := gur.db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

func (gur gormUserRepo) Update(u domain.User) error {
	user := newGormUser(u)
	result := gur.db.
		Where("id = ?", user.ID).
		Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
