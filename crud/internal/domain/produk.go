package domain

import (
	"errors"
	"time"

	"github.com/solsteace/goody/lib/oops"
)

type Produk struct {
	ID            uint
	IdToko        uint
	IdKategori    uint
	Nama          string
	Slug          string
	HargaReseller uint
	HargaKonsumen uint
	Stok          uint
	Deskripsi     string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	FotoProduk []FotoProduk
}

func NewProduk(
	id *uint,
	idToko,
	idKategori uint,
	nama,
	slug string,
	hargaReseller,
	hargaKonsumen,
	stok uint,
	deskripsi string,
	createdAt,
	updatedAt time.Time,
	fotoProduk []FotoProduk,
) (Produk, error) {
	var produkId uint = 0
	if id != nil {
		produkId = *id
	}

	switch {
	case stok < 0:
		return Produk{}, oops.BadValues{
			Err: errors.New("`stok` should be a positive number"),
			Msg: "Nilai stok harus merupakan bilangan positif"}
	case hargaKonsumen < 0:
		return Produk{}, oops.BadValues{
			Err: errors.New("`hargaKonsumen` should be a positive number"),
			Msg: "Nilai harga konsumen harus merupakan bilangan positif"}
	case hargaReseller < 0:
		return Produk{}, oops.BadValues{
			Err: errors.New("`hargaReseller` should be a positive number"),
			Msg: "Nilai harga reseller harus merupakan bilangan positif"}

	}

	for _, fp := range fotoProduk {
		if fp.IdProduk != produkId {
			return Produk{}, oops.BadValues{
				Err: errors.New("`FotoProduk` not belong to `Produk` found"),
				Msg: "id produk pada foto harus sama dengan id produk"}
		}
	}

	p := Produk{
		ID:            *id,
		IdToko:        idToko,
		IdKategori:    idKategori,
		Nama:          nama,
		Slug:          slug,
		HargaReseller: hargaReseller,
		HargaKonsumen: hargaKonsumen,
		Stok:          stok,
		Deskripsi:     deskripsi,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt}
	return p, nil
}
