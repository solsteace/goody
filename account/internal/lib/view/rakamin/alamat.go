package rakamin

import "github.com/solsteace/goody/account/internal/domain"

type Alamat struct {
	Id           uint   `json:"id"`
	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp" `
	DetailAlamat string `json:"detail_alamat"`
}

func (v viewer) Alamat(alamat domain.Alamat) any {
	return Alamat{
		Id:           alamat.ID,
		JudulAlamat:  alamat.JudulAlamat,
		NamaPenerima: alamat.NamaPenerima,
		NoTelp:       alamat.NoTelp,
		DetailAlamat: alamat.DetailAlamat}
}

func (v viewer) ManyAlamat(alamat []domain.Alamat) []any {
	view := []any{}
	for _, a := range alamat {
		viewAlamat := Alamat{
			Id:           a.ID,
			JudulAlamat:  a.JudulAlamat,
			NamaPenerima: a.NamaPenerima,
			NoTelp:       a.NoTelp,
			DetailAlamat: a.DetailAlamat}
		view = append(view, viewAlamat)
	}
	return view
}
