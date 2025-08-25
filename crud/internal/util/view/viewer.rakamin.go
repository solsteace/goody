package view

import (
	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/crud/internal/util/api"
)

type rakamin struct {
	indoApi api.Emsifa
}

func NewRakamin(indoApi api.Emsifa) rakamin {
	return rakamin{indoApi: indoApi}
}

func (v rakamin) useIndoApi(
	idProvinsi, idKota string,
) (<-chan map[string]any, <-chan map[string]any) {
	provinsi := make(chan map[string]any, 1)
	kota := make(chan map[string]any, 1)
	v.indoApi.GetProvinceAndRegencyById(idProvinsi, idKota, provinsi, kota)
	return provinsi, kota
}

func (v rakamin) Login(user domain.User, accessToken, refreshToken string) any {
	provinsi, kota := v.useIndoApi(user.IdProvinsi, user.IdKota)
	return struct {
		Nama         string         `json:"nama"`
		NoTelp       string         `json:"no_telp"`
		TanggalLahir string         `json:"tanggal_lahir"`
		Tentang      string         `json:"tentang"`
		Pekerjaan    string         `json:"pekerjaan"`
		Email        string         `json:"email"`
		IdProvinsi   map[string]any `json:"id_provinsi"`
		IdKota       map[string]any `json:"id_kota"`
		Token        string         `json:"token"`
	}{
		Nama:         user.Nama,
		NoTelp:       user.NoTelp,
		TanggalLahir: user.TanggalLahir.Format("02/01/2006"),
		Tentang:      user.Tentang,
		Pekerjaan:    user.Pekerjaan,
		Email:        user.Email,
		IdProvinsi:   <-provinsi,
		IdKota:       <-kota,
		Token:        accessToken}
}

type rakaminUser struct {
	Nama         string         `json:"nama"`
	NoTelp       string         `json:"no_telp"`
	TanggalLahir string         `json:"tanggal_lahir"`
	Tentang      string         `json:"tentang"`
	Pekerjaan    string         `json:"pekerjaan"`
	Email        string         `json:"email"`
	IdProvinsi   map[string]any `json:"id_provinsi"`
	IdKota       map[string]any `json:"id_kota"`
}

func (r rakamin) User(user domain.User) any {
	provinsi, kota := r.useIndoApi(user.IdProvinsi, user.IdKota)
	return rakaminUser{
		Nama:         user.Nama,
		NoTelp:       user.NoTelp,
		TanggalLahir: user.TanggalLahir.Format("02/01/2006"),
		Tentang:      user.Tentang,
		Pekerjaan:    user.Pekerjaan,
		Email:        user.Email,
		IdProvinsi:   <-provinsi,
		IdKota:       <-kota,
	}
}

func (r rakamin) ManyUser(user []domain.User) []any {
	view := []any{}
	for _, u := range user {
		view = append(view, r.User(u))
	}
	return view
}

type rakaminAlamat struct {
	Id           uint   `json:"id"`
	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp" `
	DetailAlamat string `json:"detail_alamat"`
}

func (r rakamin) Alamat(alamat domain.Alamat) any {
	return rakaminAlamat{
		Id:           alamat.ID,
		JudulAlamat:  alamat.JudulAlamat,
		NamaPenerima: alamat.NamaPenerima,
		NoTelp:       alamat.NoTelp,
		DetailAlamat: alamat.DetailAlamat}
}

func (r rakamin) ManyAlamat(alamat []domain.Alamat) []any {
	view := []any{}
	for _, a := range alamat {
		view = append(view, r.Alamat(a))
	}
	return view
}

type rakaminKategori struct {
	ID   uint   `json:"id"`
	Nama string `json:"nama_category"`
}

func (r rakamin) Kategori(kategori domain.Kategori) any {
	return rakaminKategori{ID: kategori.ID, Nama: kategori.Nama}
}

func (r rakamin) ManyKategori(kategori []domain.Kategori) []any {
	view := []any{}
	for _, k := range kategori {
		view = append(view, r.Kategori(k))
	}
	return view
}

type rakaminToko struct {
	Id       uint   `json:"id"`
	IdUser   uint   `json:"user_id"`
	NamaToko string `json:"nama_toko"`
	UrlFoto  string `json:"url_foto"`
}

func (r rakamin) Toko(t domain.Toko) any {
	return rakaminToko{
		Id:       t.ID,
		IdUser:   t.IdUser,
		NamaToko: t.NamaToko,
		UrlFoto:  t.UrlFoto}
}

func (r rakamin) ManyToko(toko []domain.Toko) []any {
	view := []any{}
	for _, t := range toko {
		view = append(view, r.Toko(t))
	}
	return view
}

type rakaminProduk struct {
}

func (r rakamin) Produk(produk domain.Produk) any {
	return rakaminProduk{}
}

func (r rakamin) ManyProduk(produk []domain.Produk) []any {
	view := []any{}
	for _, t := range produk {
		view = append(view, r.Produk(t))
	}
	return view
}
