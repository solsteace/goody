package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/lib/oops"
	"gorm.io/gorm"
)

type gormKategoriRow struct {
	ID        uint      `gorm:"column:id"`
	Nama      string    `gorm:"column:nama_kategori"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (row gormKategoriRow) TableName() string {
	return "kategori"
}

func newGormKategoriRow(c domain.Kategori) gormKategoriRow {
	return gormKategoriRow{
		ID:        c.ID,
		Nama:      c.Nama,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt}
}

func (row gormKategoriRow) toKategori() (domain.Kategori, error) {
	return domain.NewKategori(
		&row.ID,
		row.Nama,
		row.CreatedAt,
		row.UpdatedAt)
}

type gormKategori struct {
	db *gorm.DB
}

func NewGormKategori(db *gorm.DB) gormKategori {
	return gormKategori{db: db}
}

func (gk gormKategori) GetMany(page, limit int) ([]domain.Kategori, error) {
	rows := new([]gormKategoriRow)
	result := gk.db.
		Offset(kategoriOffset(page, limit)).
		Limit(limit).
		Find(&rows)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrRecordNotFound):
			return []domain.Kategori{}, oops.NotFound{
				Err: result.Error,
				Msg: fmt.Sprintf("Kategori tidak ditemukan")}
		default:
			return []domain.Kategori{}, result.Error
		}
	}

	kategori := []domain.Kategori{}
	for _, r := range *rows {
		k, err := r.toKategori()
		if err != nil {
			return []domain.Kategori{}, err
		}
		kategori = append(kategori, k)
	}
	return kategori, nil
}

func (gk gormKategori) GetById(id uint) (domain.Kategori, error) {
	row := new(gormKategoriRow)
	result := gk.db.
		Where("id = ?", id).
		First(&row)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrRecordNotFound):
			return domain.Kategori{}, oops.NotFound{
				Err: result.Error,
				Msg: fmt.Sprintf("Kategori(id: %d) tidak ditemukan", id)}
		default:
			return domain.Kategori{}, result.Error
		}
	}
	return row.toKategori()
}

func (gk gormKategori) Create(k domain.Kategori) (uint, error) {
	row := newGormKategoriRow(k)
	result := gk.db.Create(&row)
	if result.Error != nil {
		return 0, result.Error
	}
	return row.ID, nil
}

func (gk gormKategori) Update(k domain.Kategori) error {
	row := newGormKategoriRow(k)
	result := gk.db.
		Where("id = ?", row.ID).
		Updates(&row)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (gk gormKategori) DeleteById(id uint) error {
	row := new(gormKategoriRow)
	result := gk.db.
		Where("id = ?", id).
		Delete(&row)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
