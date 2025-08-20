package domain

import "time"

type Produk struct {
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

func NewProduk(
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
) (Produk, error) {
	// TODO: some validation
	p := Produk{
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

func (p Produk) WithId(id uint) Produk {
	p.ID = id
	return p
}
