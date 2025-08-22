package rakamin

import "github.com/solsteace/goody/account/internal/domain"

func (v viewer) Login(user domain.User, accessToken, refreshToken string) any {
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
