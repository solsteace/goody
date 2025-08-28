package domain

import (
	"errors"
	"time"

	"github.com/solsteace/goody/lib/oops"
)

type LogProduk struct {
	ID                uint
	IdDetailTransaksi uint
	IdProduk          uint
	IdToko            uint
	Nama              string
	Slug              string
	HargaReseller     uint
	HargaKonsumen     uint
	Deskripsi         string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Query
	Toko       Toko
	Kategori   Kategori
	FotoProduk []FotoProduk
}

func NewLogProduk(
	id *uint,
	idDetailTransaksi,
	idProduk,
	idToko uint,
	namaProduk,
	slug string,
	hargaReseller,
	hargaKonsumen uint,
	deskripsi string,
	createdAt,
	updatedAt time.Time,
) (LogProduk, error) {
	var logProdukId uint = 0
	if id != nil {
		logProdukId = *id
	}

	switch {
	case hargaKonsumen < 0:
		return LogProduk{}, oops.BadValues{
			Err: errors.New("`hargaKonsumen` should be a positive number"),
			Msg: "Nilai harga konsumen harus merupakan bilangan positif"}
	case hargaReseller < 0:
		return LogProduk{}, oops.BadValues{
			Err: errors.New("`hargaReseller` should be a positive number"),
			Msg: "Nilai harga reseller harus merupakan bilangan positif"}
	}

	lp := LogProduk{
		ID:                logProdukId,
		IdDetailTransaksi: idDetailTransaksi,
		IdProduk:          idProduk,
		IdToko:            idToko,
		Nama:              namaProduk,
		Slug:              slug,
		HargaReseller:     hargaReseller,
		HargaKonsumen:     hargaKonsumen,
		Deskripsi:         deskripsi,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt}
	return lp, nil
}

func (lp *LogProduk) WithKategori(k Kategori) {
	lp.Kategori = k
}

func (lp *LogProduk) WithToko(t Toko) {
	lp.Toko = t
}

func (lp *LogProduk) WithFotoProduk(fotoProduk []FotoProduk) error {
	for _, fp := range fotoProduk {
		if fp.IdProduk != lp.IdProduk {
			return oops.BadValues{
				Err: errors.New("`FotoProduk` not belong to `LogProduk` found"),
				Msg: "id produk pada foto harus sama dengan id log produk"}
		}
	}
	return nil
}
