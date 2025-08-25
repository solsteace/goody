package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/lib/oops"
	"gorm.io/gorm"
)

// Proxy object between persistence layer using Gorm and `User` domain object
type gormUserRow struct {
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
func (gu gormUserRow) TableName() string {
	return "users"
}

func (gu gormUserRow) toUser() (domain.User, error) {
	return domain.NewUser(
		&gu.ID,
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

func newGormUserRow(user domain.User) gormUserRow {
	return gormUserRow{
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
		CreatedAt:    user.CreatedAt}
}

type gormUser struct {
	db *gorm.DB
}

func NewGormUser(db *gorm.DB) gormUser {
	return gormUser{db: db}
}

func (gu gormUser) Migrate() {
	gu.db.AutoMigrate(new(gormUserRow))
}

func (gu gormUser) GetById(id uint) (domain.User, error) {
	row := new(gormUserRow)
	result := gu.db.
		Where("id = ?", id).
		First(&row)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrRecordNotFound):
			return domain.User{}, oops.NotFound{
				Err: result.Error,
				Msg: fmt.Sprintf("User(id: %d) tidak ditemukan", id)}
		default:
			return domain.User{}, result.Error
		}
	}

	return row.toUser()
}

func (gu gormUser) GetByPhoneNumber(noTelp string) (domain.User, error) {
	row := new(gormUserRow)
	result := gu.db.
		Where("no_telp = ?", noTelp).
		First(&row)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrRecordNotFound):
			return domain.User{}, oops.NotFound{
				Err: result.Error,
				Msg: fmt.Sprintf("User(telepon: %s) tidak ditemukan", noTelp)}
		default:
			return domain.User{}, result.Error
		}
	}

	return row.toUser()
}

func (gu gormUser) Create(u domain.User) (uint, error) {
	user := newGormUserRow(u)
	err := gu.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&user)
		if result.Error != nil {
			switch {
			case errors.Is(result.Error, gorm.ErrDuplicatedKey):
				return oops.BadValues{
					Err: result.Error,
					Msg: fmt.Sprintf(
						"Email `%s` atau telepon `%s` telah digunakan user lain",
						u.Email, u.NoTelp)}
			default:
				return result.Error
			}
		}

		toko := &gormTokoRow{
			IdUser:   user.ID,
			NamaToko: fmt.Sprintf("Toko %s", user.Nama),
			UrlFoto:  ""}
		result = tx.Create(&toko)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (gu gormUser) Update(u domain.User) error {
	user := newGormUserRow(u)
	result := gu.db.
		Where("id = ?", user.ID).
		Updates(user)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrDuplicatedKey):
			return oops.BadValues{
				Err: result.Error,
				Msg: fmt.Sprintf(
					"Email `%s` atau telepon `%s` telah digunakan user lain",
					u.Email, u.NoTelp)}
		default:
			return result.Error
		}
	}
	return nil
}
