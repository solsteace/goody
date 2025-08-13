package repository

import (
	"time"

	"github.com/solsteace/goody/account/internal/domain"
	"github.com/solsteace/goody/account/internal/lib/errors"
	"gorm.io/gorm"
)

// Proxy object between persistence layer using Gorm
// and `Alamat` domain object
type gormAlamatRow struct {
	ID           uint      `gorm:"column:id"`
	IdUser       uint      `gorm:"column:user_id"`
	JudulAlamat  string    `gorm:"column:judul_alamat"`
	NamaPenerima string    `gorm:"column:nama_penerima"`
	NoTelp       string    `gorm:"column:no_telp"`
	DetailAlamat string    `gorm:"column:detail_alamat"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (row gormAlamatRow) TableName() string {
	return "addresses"
}

func (row gormAlamatRow) toAlamat() (domain.Alamat, error) {
	alamat, err := domain.NewAlamat(
		row.IdUser,
		row.JudulAlamat,
		row.NamaPenerima,
		row.NoTelp,
		row.DetailAlamat,
		row.CreatedAt,
		row.UpdatedAt)
	if err != nil {
		return domain.Alamat{}, err
	}
	return alamat.WithId(row.ID), nil
}

func newGormAlamatRow(alamat domain.Alamat) gormAlamatRow {
	return gormAlamatRow{
		ID:           alamat.ID,
		IdUser:       alamat.UserId,
		JudulAlamat:  alamat.JudulAlamat,
		NamaPenerima: alamat.NamaPenerima,
		NoTelp:       alamat.NoTelp,
		DetailAlamat: alamat.DetailAlamat,
		CreatedAt:    alamat.CreatedAt,
		UpdatedAt:    alamat.UpdatedAt}
}

type gormAlamat struct {
	db *gorm.DB
}

func NewGormAlamat(db *gorm.DB) gormAlamat {
	return gormAlamat{db: db}
}

func (ga gormAlamat) GetManyByUserId(userId, offset, limit uint) ([]domain.Alamat, error) {
	rows := new([]gormAlamatRow)
	result := ga.db.
		Where("user_id = ?", userId).
		Find(&rows)
	if result.Error != nil {
		return []domain.Alamat{}, errors.Standardize(result.Error)
	}

	daftarAlamat := []domain.Alamat{}
	for _, r := range *rows {
		alamat, err := r.toAlamat()
		if err != nil {
			return []domain.Alamat{}, err
		}
		daftarAlamat = append(daftarAlamat, alamat)
	}
	return daftarAlamat, nil
}

func (ga gormAlamat) GetById(id uint) (domain.Alamat, error) {
	row := new(gormAlamatRow)
	result := ga.db.
		Where("id = ?", id).
		First(&row)
	if result.Error != nil {
		return domain.Alamat{}, errors.Standardize(result.Error)
	}

	alamat, err := row.toAlamat()
	if err != nil {
		return domain.Alamat{}, errors.Standardize(result.Error)
	}
	return alamat, errors.Standardize(result.Error)
}

func (ga gormAlamat) Update(alamat domain.Alamat) error {
	row := newGormAlamatRow(alamat)
	result := ga.db.
		Where("id = ?", row.ID).
		Updates(row)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ga gormAlamat) DeleteById(id uint) error {
	row := new(gormAlamatRow)
	result := ga.db.
		Where("id = ?", id).
		Delete(&row)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ga gormAlamat) Create(alamat domain.Alamat) (uint, error) {
	row := newGormAlamatRow(alamat)
	result := ga.db.Create(&row)
	if result.Error != nil {
		return 0, result.Error
	}
	return row.ID, nil
}
