package rakamin

import "github.com/solsteace/goody/account/internal/domain"

type User struct {
	Nama         string         `json:"nama"`
	NoTelp       string         `json:"no_telp"`
	TanggalLahir string         `json:"tanggal_lahir"`
	Tentang      string         `json:"tentang"`
	Pekerjaan    string         `json:"pekerjaan"`
	Email        string         `json:"email"`
	IdProvinsi   map[string]any `json:"id_provinsi"`
	IdKota       map[string]any `json:"id_kota"`
}

func (v viewer) User(user domain.User) any {
	provinsi, kota := v.useIndoApi(user.IdProvinsi, user.IdKota)
	return User{
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

func (v viewer) ManyUser(user []domain.User) []any {
	view := []any{}
	for _, u := range user {
		provinsi, kota := v.useIndoApi(u.IdProvinsi, u.IdKota)
		viewUser := User{
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
