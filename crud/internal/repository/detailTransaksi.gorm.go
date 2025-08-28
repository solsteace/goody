package repository

import (
	"time"

	"github.com/solsteace/goody/crud/internal/domain"
)

type gormDetailTransaksiRow struct {
	ID          uint             `gorm:"column:id"`
	IdTransaksi uint             `gorm:"column:id_transaksi"`
	IdToko      uint             `gorm:"column:id_toko"`
	Kuantitas   uint             `gorm:"column:kuantitas"`
	HargaTotal  uint             `gorm:"column:harga_total"`
	CreatedAt   time.Time        `gorm:"column:created_at"`
	UpdatedAt   time.Time        `gorm:"column:updated_at"`
	LogProduk   gormLogProdukRow `gorm:"foreignKey:IdDetailTransaksi"`

	// Query
	Toko gormTokoRow `gorm:"foreignKey:IdToko"`
}

func (row gormDetailTransaksiRow) TableName() string {
	return "detail_transaksi"
}

func (row gormDetailTransaksiRow) toDetailTransaksi() (domain.DetailTransaksi, error) {
	logProduk, err := row.LogProduk.toLogProduk()
	if err != nil {
		return domain.DetailTransaksi{}, err
	}

	detailTransaksi, err := domain.NewDetailTransaksi(
		&row.ID,
		&row.IdTransaksi,
		row.IdToko,
		row.Kuantitas,
		row.HargaTotal,
		row.CreatedAt,
		row.UpdatedAt,
		logProduk)
	if err != nil {
		return domain.DetailTransaksi{}, err
	}

	toko, err := row.Toko.toToko()
	if err != nil {
		return domain.DetailTransaksi{}, err
	}

	detailTransaksi.WithToko(toko)
	return detailTransaksi, nil
}

func newGormDetailTransaksiRow(dt domain.DetailTransaksi) gormDetailTransaksiRow {
	return gormDetailTransaksiRow{
		ID:          dt.ID,
		IdTransaksi: dt.IdTransaksi,
		IdToko:      dt.IdToko,
		Kuantitas:   dt.Kuantitas,
		HargaTotal:  dt.HargaTotal,
		CreatedAt:   dt.CreatedAt,
		UpdatedAt:   dt.UpdatedAt,
		LogProduk:   newGormLogProdukRow(dt.LogProduk)}
}
