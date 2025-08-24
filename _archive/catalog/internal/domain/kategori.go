package domain

import "time"

type Kategori struct {
	ID        uint      `json:"id"`
	Nama      string    `json:"nama_kategori"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewKategori(
	id *uint,
	nama string,
	createdAt time.Time,
	updatedAt time.Time,
) (Kategori, error) {
	// TODO: domain validation
	kategori := Kategori{
		ID:        *id,
		Nama:      nama,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt}
	return kategori, nil
}
