package view

import "github.com/solsteace/goody/catalog/internal/domain"

type rakamin struct{}

func NewRakamin() rakamin {
	return rakamin{}
}

type rakaminKategori struct {
	ID   uint   `json:"id"`
	Nama string `json:"nama_category"`
}

func (r rakamin) Kategori(k domain.Kategori) any {
	return rakaminKategori{ID: k.ID, Nama: k.Nama}
}

func (r rakamin) ManyKategori(k []domain.Kategori) []any {
	view := []any{}
	for _, kv := range k {
		viewKategori := rakaminKategori{ID: kv.ID, Nama: kv.Nama}
		view = append(view, viewKategori)
	}
	return view
}

type toko struct {
	Id       uint   `json:"id"`
	IdUser   uint   `json:"user_id"`
	NamaToko string `json:"nama_toko"`
	UrlFoto  string `json:"url_foto"`
}

func (r rakamin) Toko(t domain.Toko) any {
	return toko{
		Id:       t.ID,
		IdUser:   t.IdUser,
		NamaToko: t.NamaToko,
		UrlFoto:  t.UrlFoto}
}

func (r rakamin) ManyToko(t []domain.Toko) []any {
	view := []any{}
	for _, tv := range t {
		viewToko := toko{
			Id:       tv.ID,
			IdUser:   tv.IdUser,
			NamaToko: tv.NamaToko,
			UrlFoto:  tv.UrlFoto}
		view = append(view, viewToko)
	}
	return view
}
