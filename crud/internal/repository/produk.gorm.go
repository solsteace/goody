package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/lib/oops"
	"gorm.io/gorm"
)

type gormProdukRow struct {
	ID            uint      `gorm:"column:id;primaryKey"`
	IdToko        uint      `gorm:"column:id_toko"`
	IdKategori    uint      `gorm:"column:id_kategori"`
	Nama          string    `gorm:"column:nama_produk"`
	Slug          string    `gorm:"column:slug"`
	HargaReseller uint      `gorm:"column:harga_reseller"`
	HargaKonsumen uint      `gorm:"column:harga_konsumen"`
	Stok          uint      `gorm:"column:stok"`
	Deskripsi     string    `gorm:"column:deskripsi"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`

	FotoProduk []gormFotoProdukRow `gorm:"foreignKey:IdProduk"`
	Toko       gormTokoRow         `gorm:"foreignKey:IdToko"`
	Kategori   gormKategoriRow     `gorm:"foreignKey:IdKategori"`
}

func (row gormProdukRow) TableName() string {
	return "produk"
}

func NewGormProdukRow(p domain.Produk) gormProdukRow {
	fotoProdukRows := []gormFotoProdukRow{}
	for _, fp := range p.FotoProduk {
		fotoProdukRows = append(fotoProdukRows, NewGormFotoProdukRow(fp))
	}

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
		UpdatedAt:     p.UpdatedAt,
		FotoProduk:    fotoProdukRows}
}

func (row gormProdukRow) toProduk() (domain.Produk, error) {
	fotoProduk := []domain.FotoProduk{}
	for _, row := range row.FotoProduk {
		fp, err := domain.NewFotoProduk(
			&row.ID,
			&row.IdProduk,
			row.Url,
			row.CreatedAt,
			row.UpdatedAt)
		if err != nil {
			return domain.Produk{}, err
		}
		fotoProduk = append(fotoProduk, fp)
	}

	kategori, err := row.Kategori.toKategori()
	if err != nil {
		return domain.Produk{}, err
	}

	toko, err := row.Toko.toToko()
	if err != nil {
		return domain.Produk{}, err
	}

	produk, err := domain.NewProduk(
		&row.ID,
		&row.IdToko,
		&row.IdKategori,
		row.Nama,
		row.Slug,
		int(row.HargaReseller),
		int(row.HargaKonsumen),
		int(row.Stok),
		row.Deskripsi,
		row.CreatedAt,
		row.UpdatedAt,
		fotoProduk)
	if err != nil {
		return domain.Produk{}, err
	}

	produk.WithToko(toko)
	produk.WithKategori(kategori)
	return produk, nil
}

type gormProduk struct {
	db *gorm.DB
}

func NewGormProduk(db *gorm.DB) gormProduk {
	return gormProduk{db: db}
}

func (gp gormProduk) GetMany(params produkQueryParams) ([]domain.Produk, error) {
	query := gp.db.
		Preload("Toko").
		Preload("FotoProduk").
		Preload("Kategori").
		Where("harga_reseller <= ?", params.maxHarga).
		Where("harga_konsumen <= ?", params.maxHarga).
		Where("harga_reseller >= ?", params.minHarga).
		Where("harga_konsumen >= ?", params.minHarga).
		Where("nama_produk LIKE ?", "%"+params.nama+"%").
		Offset(params.offset()).
		Limit(params.limit)
	if params.kategoriId != 0 {
		query = query.Where("id_kategori = ?", params.kategoriId)
	}
	if params.tokoId != 0 {
		query = query.Where("id_toko = ?", params.tokoId)
	}

	rows := new([]gormProdukRow)
	result := query.Find(rows)
	if result.Error != nil {
		return []domain.Produk{}, result.Error
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
		Preload("Toko").
		Preload("FotoProduk").
		Preload("Kategori").
		Where("id = ?", id).
		First(row)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrRecordNotFound):
			return domain.Produk{}, oops.NotFound{
				Err: result.Error,
				Msg: fmt.Sprintf("Produk(id: %d) tidak ditemukan", id)}
		default:
			return domain.Produk{}, result.Error
		}
	}

	produk, err := row.toProduk()
	if err != nil {
		return domain.Produk{}, err
	}
	return produk, nil
}

func (gp gormProduk) Create(p domain.Produk) (uint, error) {
	row := NewGormProdukRow(p)
	result := gp.db.
		Omit("Toko", "Category"). // No need these to be touched
		Create(&row)
	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrForeignKeyViolated):
			return 0, oops.BadValues{
				Err: result.Error,
				Msg: fmt.Sprintf("kategori(id:%d) atau toko(id:%d) tidak ditemukan",
					p.IdKategori, p.IdToko),
			}
		}
		return 0, result.Error
	}
	return row.ID, nil
}

func (gp gormProduk) Update(p domain.Produk) error {
	err := gp.db.Transaction(func(tx *gorm.DB) error {
		fotoProdukRow := new(gormFotoProdukRow)
		result := tx.
			Where("id_produk = ?", p.ID).
			Delete(fotoProdukRow)
		if result.Error != nil {
			return result.Error
		}

		row := NewGormProdukRow(p)
		result = tx.
			Omit("Toko", "Category"). // No need these to be touched
			Where("id = ?", p.ID).
			Updates(&row)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrForeignKeyViolated):
			return oops.BadValues{
				Err: err,
				Msg: fmt.Sprintf("kategori(id:%d) atau toko(id:%d) tidak ditemukan",
					p.IdKategori, p.IdToko),
			}
		}
		return err
	}
	return nil
}

func (gp gormProduk) DeleteById(id uint) error {
	err := gp.db.Transaction(func(tx *gorm.DB) error {
		fotoProdukRow := new(gormFotoProdukRow)
		result := tx.
			Where("id_produk = ?", id).
			Delete(fotoProdukRow)
		if result.Error != nil {
			return result.Error
		}

		produkRow := new(gormProdukRow)
		result = tx.
			Where("id = ?", id).
			Delete(&produkRow)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
