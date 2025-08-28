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

type rakaminFotoProduk struct {
	Id       uint   `json:"id"`
	IdProduk uint   `json:"product_id"`
	Url      string `json:"url"`
}

func (r rakamin) FotoProduk(fp domain.FotoProduk) any {
	return rakaminFotoProduk{
		Id:       fp.ID,
		IdProduk: fp.IdProduk,
		Url:      fp.Url}
}

func (r rakamin) ManyFotoProduk(fotoProduk []domain.FotoProduk) []any {
	view := []any{}
	for _, fp := range fotoProduk {
		view = append(view, r.FotoProduk(fp))
	}
	return view
}

type rakaminProduk struct {
	Id            uint   `json:"id"`
	Nama          string `json:"nama_produk"`
	Slug          string `json:"slug"`
	HargaReseller uint   `json:"harga_reseller"`
	HargaKonsumen uint   `json:"harga_konsumen"`
	Stok          uint   `json:"stok"`
	Deskripsi     string `json:"deskripsi"`

	Toko       rakaminToko         `json:"toko"`
	Kategori   rakaminKategori     `json:"category"`
	FotoProduk []rakaminFotoProduk `json:"photos"`
}

func (r rakamin) Produk(produk domain.Produk) any {
	toko, _ := r.Toko(produk.Toko).(rakaminToko)
	kategori, _ := r.Kategori(produk.Kategori).(rakaminKategori)
	fotoProduk := []rakaminFotoProduk{}
	for _, fp := range produk.FotoProduk {
		fpView, _ := r.FotoProduk(fp).(rakaminFotoProduk)
		fotoProduk = append(fotoProduk, fpView)
	}

	return rakaminProduk{
		Id:            produk.ID,
		Nama:          produk.Nama,
		Slug:          produk.Slug,
		HargaReseller: produk.HargaReseller,
		HargaKonsumen: produk.HargaKonsumen,
		Stok:          produk.Stok,
		Deskripsi:     produk.Deskripsi,
		Toko:          toko,
		Kategori:      kategori,
		FotoProduk:    fotoProduk}
}

func (r rakamin) ManyProduk(produk []domain.Produk) []any {
	view := []any{}
	for _, t := range produk {
		view = append(view, r.Produk(t))
	}
	return view
}

type rakaminLogProduk struct {
	ID            uint   `json:"id"`
	Nama          string `json:"nama_produk"`
	Slug          string `json:"slug"`
	HargaReseller uint   `json:"harga_reseller"`
	HargaKonsumen uint   `json:"harga_konsumen"`
	Deskripsi     string `json:"deskripsi"`

	Toko     rakaminToko     `json:"toko"`
	Kategori rakaminKategori `json:"category"`
}

func (r rakamin) LogProduk(lp domain.LogProduk) any {
	toko := r.Toko(lp.Toko).(rakaminToko)
	Kategori := r.Kategori(lp.Kategori).(rakaminKategori)

	return rakaminLogProduk{
		ID:            lp.ID,
		Nama:          lp.Nama,
		Slug:          lp.Slug,
		HargaReseller: lp.HargaReseller,
		HargaKonsumen: lp.HargaKonsumen,
		Deskripsi:     lp.Deskripsi,

		Toko:     toko,
		Kategori: Kategori}
}

func (r rakamin) ManyLogProduk(logProduk []domain.LogProduk) []any {
	view := []any{}
	for _, lp := range logProduk {
		view = append(view, r.LogProduk(lp))
	}
	return view
}

type rakaminDetailTransaksi struct {
	Kuantitas  uint             `json:"kuantitas"`
	HargaTotal uint             `json:"harga_total"`
	LogProduk  rakaminLogProduk `json:"product"`
}

func (r rakamin) DetailTransaksi(dt domain.DetailTransaksi) any {
	logProduk := r.LogProduk(dt.LogProduk).(rakaminLogProduk)
	return rakaminDetailTransaksi{
		Kuantitas:  dt.Kuantitas,
		HargaTotal: dt.HargaTotal,
		LogProduk:  logProduk}
}

func (r rakamin) ManyDetailTransaksi(detailTx []domain.DetailTransaksi) []any {
	view := []any{}
	for _, dt := range detailTx {
		view = append(view, r.DetailTransaksi(dt))
	}
	return view
}

type rakaminTransaksi struct {
	Id          uint   `json:"id"`
	HargaTotal  uint   `json:"harga_total"`
	KodeInvoice string `json:"kode_invoice"`
	MetodeBayar string `json:"metode_bayar"`

	Alamat          rakaminAlamat            `json:"alamat_kirim"`
	DetailTransaksi []rakaminDetailTransaksi `json:"detail_trx"`
}

func (r rakamin) Transaksi(t domain.Transaksi) any {
	alamat, _ := r.Alamat(t.Alamat).(rakaminAlamat)
	detailTransaksi := []rakaminDetailTransaksi{}
	for _, dt := range t.DetailTransaksi {
		dtView := r.DetailTransaksi(dt).(rakaminDetailTransaksi)
		detailTransaksi = append(detailTransaksi, dtView)
	}

	return rakaminTransaksi{
		Id:              t.ID,
		HargaTotal:      t.HargaTotal,
		MetodeBayar:     t.MetodeBayar,
		KodeInvoice:     t.KodeInvoice,
		Alamat:          alamat,
		DetailTransaksi: detailTransaksi}
}

func (r rakamin) ManyTransaksi(transaksi []domain.Transaksi) []any {
	view := []any{}
	for _, t := range transaksi {
		view = append(view, r.Transaksi(t))
	}
	return view
}
