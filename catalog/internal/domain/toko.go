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
	idUser uint,
	namaToko,
	urlFoto string,
	createdAt time.Time,
	updatedAt time.Time,
) (Toko, error) {
	// TODO: Domain validation...
	t := Toko{
		IdUser:    idUser,
		NamaToko:  namaToko,
		UrlFoto:   urlFoto,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt}
	return t, nil
}

func (t Toko) WithId(id uint) Toko {
	t.ID = id
	return t
}
