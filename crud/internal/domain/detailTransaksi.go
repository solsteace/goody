package domain

import (
	"errors"
	"time"

	"github.com/solsteace/goody/lib/oops"
)

type DetailTransaksi struct {
	ID          uint
	IdTransaksi uint
	IdToko      uint
	Kuantitas   uint
	HargaTotal  uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LogProduk   LogProduk

	// Query
	Toko Toko
}

func NewDetailTransaksi(
	id,
	idTransaksi *uint,
	idToko,
	kuantitas,
	hargaTotal uint,
	createdAt,
	updatedAt time.Time,
	logProduk LogProduk,
) (DetailTransaksi, error) {
	var actualIdDetailTransaksi uint = 0
	if id != nil {
		actualIdDetailTransaksi = *id
	}

	var actualIdTransaksi uint = 0
	if idTransaksi != nil {
		actualIdTransaksi = *idTransaksi
	}

	if logProduk.IdDetailTransaksi != actualIdDetailTransaksi {
		return DetailTransaksi{}, oops.BadValues{
			Err: errors.New("`LogProduk` not belong to `DetailTransaksi` found")}
	}

	// TODO: domain validation
	dt := DetailTransaksi{
		ID:          actualIdDetailTransaksi,
		IdTransaksi: actualIdTransaksi,
		IdToko:      idToko,
		Kuantitas:   kuantitas,
		HargaTotal:  hargaTotal,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		LogProduk:   logProduk}
	return dt, nil
}

func (dt *DetailTransaksi) WithToko(t Toko) {
	dt.Toko = t
}
