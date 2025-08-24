package domain

import (
	"errors"
	"time"

	"github.com/solsteace/goody/catalog/internal/domain/entity"
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

	FotoProduk []entity.FotoProduk
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
	fotoProduk []entity.FotoProduk,
) (Produk, error) {
	var produkId uint = 0
	if id != nil {
		produkId = *id
	}

	// TODO: some validation
	for _, fp := range fotoProduk {
		if fp.IdProduk != produkId {
			return Produk{}, errors.New("`FotoProduk` not belong to `Produk` found")
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
