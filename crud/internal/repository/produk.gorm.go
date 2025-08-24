package repository

import (
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
	"gorm.io/gorm"
)

type gormProdukRow struct {
	ID            uint      `gorm:"id"`
	IdToko        uint      `gorm:"id_toko"`
	IdKategori    uint      `gorm:"id_kategori"`
	Nama          string    `gorm:"nama_produk"`
	Slug          string    `gorm:"slug"`
	HargaReseller uint      `gorm:"harga_reseller"`
	HargaKonsumen uint      `gorm:"harga_konsumen"`
	Stok          uint      `gorm:"stok"`
	Deskripsi     string    `gorm:"deskripsi"`
	CreatedAt     time.Time `gorm:"created_at"`
	UpdatedAt     time.Time `gorm:"updated_at"`

	FotoProduk gormFotoProdukRow `gorm:"references:IdProduk"`
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
	return domain.NewProduk(
		&row.ID,
		row.IdToko,
		row.IdKategori,
		row.Nama,
		row.Slug,
		row.HargaReseller,
		row.HargaKonsumen,
		row.Stok,
		row.Deskripsi,
		row.CreatedAt,
		row.UpdatedAt,
		[]domain.FotoProduk{})
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
		Preload(gormProdukRow{}.TableName()).
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
		Preload(gormProdukRow{}.TableName()).
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
