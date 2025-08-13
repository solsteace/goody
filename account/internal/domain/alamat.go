package domain

import "time"

type Alamat struct {
	ID           uint      `json:"id"`
	UserId       uint      `json:"user_id"`
	JudulAlamat  string    `json:"judul_alamat"`
	NamaPenerima string    `json:"nama_penerima"`
	NoTelp       string    `json:"no_telp"`
	DetailAlamat string    `json:"detail_alamat"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewAlamat(
	userId uint,
	judulAlamat,
	namaPenerima,
	noTelp,
	detailAlamat string,
	createdAt time.Time,
	updatedAt time.Time,
) (Alamat, error) {
	// TODO: domain validation...

	a := Alamat{
		UserId:       userId,
		JudulAlamat:  judulAlamat,
		NamaPenerima: namaPenerima,
		NoTelp:       noTelp,
		DetailAlamat: detailAlamat,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt}
	return a, nil
}

func (a Alamat) WithId(id uint) Alamat {
	a.ID = id
	return a
}
