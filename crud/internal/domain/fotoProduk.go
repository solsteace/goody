package domain

import "time"

type FotoProduk struct {
	ID        uint      // `json:"id"`
	IdProduk  uint      // `json:"id_produk"`
	Url       string    // `json:"url"`
	CreatedAt time.Time // `json:"created_at"`
	UpdatedAt time.Time // `json:"updated_at"`
}

func NewFotoProduk(
	id *uint,
	idProduk *uint,
	url string,
	createdAt,
	updatedAt time.Time,
) (FotoProduk, error) {
	var idFotoProduk uint = 0
	if id != nil {
		idFotoProduk = *id
	}

	var actualIdProduk uint = 0
	if idProduk != nil {
		actualIdProduk = uint(*idProduk)
	}

	// TODO: validation
	fp := FotoProduk{
		ID:        idFotoProduk,
		IdProduk:  actualIdProduk,
		Url:       url,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt}
	return fp, nil
}
