package repository

import (
	"time"

	"github.com/solsteace/goody/catalog/internal/domain"
)

type GormCategoryRow struct {
	ID        uint      `json:"id"`
	Nama      string    `json:"nama_kategori"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (row GormCategoryRow) TableName() string {
	return "kategori"
}

func NewGormCategoryRow(c domain.Kategori) GormCategoryRow {
	return GormCategoryRow{
		ID:        c.ID,
		Nama:      c.Nama,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt}
}

func (row GormCategoryRow) toCategory() (domain.Kategori, error) {
	k, err := domain.NewKategori(
		row.Nama,
		row.CreatedAt,
		row.UpdatedAt)
	if err != nil {
		return domain.Kategori{}, err
	}
	return k.WithId(row.ID), nil
}
