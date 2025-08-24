package repository

import (
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
)

type gormFotoProdukRow struct {
	ID        uint      `gorm:"id"`
	IdProduk  uint      `gorm:"id_produk"`
	Url       string    `gorm:"url"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (row gormFotoProdukRow) TableName() string {
	return "foto_produk"
}

func NewGormFotoProdukRow(fp domain.FotoProduk) gormFotoProdukRow {
	return gormFotoProdukRow{
		ID:        fp.ID,
		IdProduk:  fp.IdProduk,
		Url:       fp.Url,
		CreatedAt: fp.CreatedAt,
		UpdatedAt: fp.UpdatedAt}
}

func (row gormFotoProdukRow) toFotoProduk() (domain.FotoProduk, error) {
	return domain.NewFotoProduk(
		&row.ID,
		row.IdProduk,
		row.Url,
		row.CreatedAt,
		row.UpdatedAt)
}
