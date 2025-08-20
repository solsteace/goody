package domain

import "time"

type Kategori struct {
	ID        uint      `json:"id"`
	Nama      string    `json:"nama_kategori"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewKategori(
	nama string,
	createdAt time.Time,
	updatedAt time.Time,
) (Kategori, error) {
	// TODO: domain validation
	kategori := Kategori{
		Nama:      nama,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt}
	return kategori, nil
}

func (k Kategori) WithId(id uint) Kategori {
	k.ID = id
	return k
}
