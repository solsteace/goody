package view

import (
	"github.com/solsteace/goody/account/internal/domain"
	"github.com/solsteace/goody/account/internal/lib/api"
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

func (v rakamin) User(user domain.User) any {
	provinsi, kota := v.useIndoApi(user.IdProvinsi, user.IdKota)
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

func (v rakamin) ManyUser(user []domain.User) []any {
	view := []any{}
	for _, u := range user {
		provinsi, kota := v.useIndoApi(u.IdProvinsi, u.IdKota)
		viewUser := rakaminUser{
			Nama:         u.Nama,
			NoTelp:       u.NoTelp,
			TanggalLahir: u.TanggalLahir.Format("02/01/2006"),
			Tentang:      u.Tentang,
			Pekerjaan:    u.Pekerjaan,
			Email:        u.Email,
			IdProvinsi:   <-provinsi,
			IdKota:       <-kota}
		view = append(view, viewUser)
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

func (v rakamin) Alamat(alamat domain.Alamat) any {
	return rakaminAlamat{
		Id:           alamat.ID,
		JudulAlamat:  alamat.JudulAlamat,
		NamaPenerima: alamat.NamaPenerima,
		NoTelp:       alamat.NoTelp,
		DetailAlamat: alamat.DetailAlamat}
}

func (v rakamin) ManyAlamat(alamat []domain.Alamat) []any {
	view := []any{}
	for _, a := range alamat {
		viewAlamat := rakaminAlamat{
			Id:           a.ID,
			JudulAlamat:  a.JudulAlamat,
			NamaPenerima: a.NamaPenerima,
			NoTelp:       a.NoTelp,
			DetailAlamat: a.DetailAlamat}
		view = append(view, viewAlamat)
	}
	return view
}
