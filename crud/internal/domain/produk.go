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
	FotoProduk    []FotoProduk

	Toko     Toko
	Kategori Kategori
}

func NewProduk(
	id,
	idToko,
	idKategori *uint,
	nama,
	slug string,
	hargaReseller,
	hargaKonsumen,
	stok int,
	deskripsi string,
	createdAt,
	updatedAt time.Time,
	fotoProduk []FotoProduk,
) (Produk, error) {
	var actualIdProduk uint = 0
	if id != nil {
		actualIdProduk = *id
	}
	var actualIdToko uint = 0
	if idToko != nil {
		actualIdToko = *idToko
	}
	var actualIdKategori uint = 0
	if idKategori != nil {
		actualIdKategori = *idKategori
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
		if fp.IdProduk != actualIdProduk {
			return Produk{}, oops.BadValues{
				Err: errors.New("`FotoProduk` not belong to `Produk` found"),
				Msg: "id produk pada foto harus sama dengan id produk"}
		}
	}

	p := Produk{
		ID:            actualIdProduk,
		IdToko:        actualIdToko,
		IdKategori:    actualIdKategori,
		Nama:          nama,
		Slug:          slug,
		HargaReseller: uint(hargaReseller),
		HargaKonsumen: uint(hargaKonsumen),
		Stok:          uint(stok),
		Deskripsi:     deskripsi,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		FotoProduk:    fotoProduk}
	return p, nil
}

func (p *Produk) WithToko(t Toko) {
	p.Toko = t
}

func (p *Produk) WithKategori(k Kategori) {
	p.Kategori = k
}
