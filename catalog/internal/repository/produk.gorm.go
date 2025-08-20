package repository

import (
	"time"

	"github.com/solsteace/goody/catalog/internal/domain"
	"gorm.io/gorm"
)

type gormProdukRow struct {
	ID            uint      `json:"id"`
	IdToko        uint      `json:"id_toko"`
	IdKategori    uint      `json:"id_kategori"`
	Nama          string    `json:"nama_produk"`
	Slug          string    `json:"slug"`
	HargaReseller uint      `json:"harga_reseller"`
	HargaKonsumen uint      `json:"harga_konsumen"`
	Stok          uint      `json:"stok"`
	Deskripsi     string    `json:"deskripsi"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (row gormProdukRow) TableName() string {
	return "produk"
}

func NewGormProdukRow(p domain.Produk) gormProdukRow {
	return gormProdukRow{
		ID:            p.ID,
		IdToko:        p.IdToko,
		IdKategori:    p.IdKategori,
		Nama:          p.Nama,
		Slug:          p.Slug,
		HargaReseller: p.HargaReseller,
		HargaKonsumen: p.HargaKonsumen,
		Stok:          p.Stok,
		Deskripsi:     p.Deskripsi,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt}
}

func (row gormProdukRow) toProduk() (domain.Produk, error) {
	p, err := domain.NewProduk(
		row.IdToko,
		row.IdKategori,
		row.Nama,
		row.Slug,
		row.HargaReseller,
		row.HargaKonsumen,
		row.Stok,
		row.Deskripsi,
		row.CreatedAt,
		row.UpdatedAt)
	if err != nil {
		return domain.Produk{}, err
	}
	return p.WithId(row.ID), nil
}

type gormProduk struct {
	db *gorm.DB
}

func NewGormProduk(db *gorm.DB) gormProduk {
	return gormProduk{db: db}
}

func (gp gormProduk) GetMany(page int, limit int) ([]domain.Produk, error) {
	rows := new([]gormProdukRow)
	result := gp.db.
		Offset(produkOffset(page, limit)).
		Limit(limit).
		Find(&rows)
	if result.Error != nil {
		return []domain.Produk{}, nil
	}

	produk := []domain.Produk{}
	for _, r := range *rows {
		p, err := r.toProduk()
		if err != nil {
			return []domain.Produk{}, err
		}
		produk = append(produk, p)
	}
	return produk, nil
}

func (gp gormProduk) GetById(id uint) (domain.Produk, error) {
	row := new(gormProdukRow)
	result := gp.db.
		Where("id = ?", id).
		First(&row)
	if result.Error != nil {
		return domain.Produk{}, nil
	}

	produk, err := row.toProduk()
	if err != nil {
		return domain.Produk{}, nil
	}
	return produk, nil
}

func (gp gormProduk) Create(p domain.Produk) (uint, error) {
	row := NewGormProdukRow(p)
	result := gp.db.Create(&row)
	if result.Error != nil {
		return 0, result.Error
	}
	return row.ID, nil
}

func (gp gormProduk) Update(p domain.Produk) error {
	row := NewGormProdukRow(p)
	result := gp.db.
		Where("id = ?", p.ID).
		Updates(&row)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (gp gormProduk) DeleteById(id uint) error {
	row := new(gormProdukRow)
	result := gp.db.
		Where("id = ?", id).
		Delete(&row)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
