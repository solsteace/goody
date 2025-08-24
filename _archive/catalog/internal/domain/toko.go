package domain

import "time"

type Toko struct {
	ID        uint      `json:"id"`
	IdUser    uint      `json:"id_user"`
	NamaToko  string    `json:"nama_toko"`
	UrlFoto   string    `json:"url_foto"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewToko(
	id *uint,
	idUser uint,
	namaToko,
	urlFoto string,
	createdAt time.Time,
	updatedAt time.Time,
) (Toko, error) {
	var tokoId uint = 0
	if id != nil {
		tokoId = *id
	}

	// TODO: Domain validation...
	t := Toko{
		ID:        tokoId,
		IdUser:    idUser,
		NamaToko:  namaToko,
		UrlFoto:   urlFoto,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt}
	return t, nil
}
