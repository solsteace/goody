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
	var kategoriId uint = 0
	if id != nil {
		kategoriId = *id
	}

	// TODO: domain validation
	kategori := Kategori{
		ID:        kategoriId,
		Nama:      nama,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt}
	return kategori, nil
}
