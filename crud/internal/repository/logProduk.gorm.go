package repository

import (
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
)

type gormLogProdukRow struct {
	ID                uint      `gorm:"column:id;primaryKey"`
	IdDetailTransaksi uint      `gorm:"column:id_detail_transaksi"`
	IdProduk          uint      `gorm:"column:id_produk"`
	IdToko            uint      `gorm:"column:id_toko"`
	Nama              string    `gorm:"column:nama_produk"`
	Slug              string    `gorm:"column:slug"`
	HargaReseller     uint      `gorm:"column:harga_reseller"`
	HargaKonsumen     uint      `gorm:"column:harga_konsumen"`
	Deskripsi         string    `gorm:"column:deskripsi"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at"`

	// Query
	Produk gormProdukRow `gorm:"foreignKey:IdProduk"`
}

func (row gormLogProdukRow) TableName() string {
	return "log_produk"
}

func (row gormLogProdukRow) toLogProduk() (domain.LogProduk, error) {
	logProduk, err := domain.NewLogProduk(
		&row.ID,
		row.IdDetailTransaksi,
		row.IdProduk,
		row.IdToko,
		row.Nama,
		row.Slug,
		row.HargaReseller,
		row.HargaKonsumen,
		row.Deskripsi,
		row.CreatedAt,
		row.UpdatedAt)
	if err != nil {
		return domain.LogProduk{}, err
	}

	produk, err := row.Produk.toProduk()
	if err != nil {
		return domain.LogProduk{}, err
	}

	if err := logProduk.WithFotoProduk(produk.FotoProduk); err != nil {
		return domain.LogProduk{}, err
	}
	logProduk.WithKategori(produk.Kategori)
	logProduk.WithToko(produk.Toko)
	return logProduk, nil
}

func newGormLogProdukRow(lp domain.LogProduk) gormLogProdukRow {
	return gormLogProdukRow{
		ID:                lp.ID,
		IdDetailTransaksi: lp.IdDetailTransaksi,
		IdProduk:          lp.IdProduk,
		IdToko:            lp.IdToko,
		Nama:              lp.Nama,
		Slug:              lp.Slug,
		HargaReseller:     lp.HargaReseller,
		HargaKonsumen:     lp.HargaKonsumen,
		Deskripsi:         lp.Deskripsi,
		CreatedAt:         lp.CreatedAt,
		UpdatedAt:         lp.UpdatedAt}
}
