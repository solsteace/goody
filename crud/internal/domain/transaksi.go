package domain

import (
	"errors"
	"time"

	"github.com/solsteace/goody/lib/oops"
)

type Transaksi struct {
	ID              uint
	IdUser          uint
	IdAlamat        uint
	HargaTotal      uint
	KodeInvoice     string
	MetodeBayar     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DetailTransaksi []DetailTransaksi

	// Query
	Alamat Alamat
}

func NewTransaksi(
	id *uint,
	idUser,
	idAlamat uint,
	hargaTotal uint,
	kodeInvoice,
	metodeBayar string,
	createdAt,
	updatedAt time.Time,
	detailTransaksi []DetailTransaksi,
) (Transaksi, error) {
	var transaksiId uint = 0
	if id != nil {
		transaksiId = *id
	}

	for _, dt := range detailTransaksi {
		if dt.IdTransaksi != transaksiId {
			return Transaksi{}, oops.BadValues{
				Err: errors.New("`DetailTransaksi` not belong to `Transaksi` found")}
		}
	}
	if hargaTotal < 0 {
		return Transaksi{}, oops.BadValues{
			Err: errors.New("`HargaTotal` must be a positive number"),
			Msg: "`HargaTotal` harus bernilai positif"}
	}

	if len(kodeInvoice) == 0 {
		return Transaksi{}, oops.BadValues{
			Err: errors.New("`KodeInvoice` cannot be an empty string"),
			Msg: "`KodeInvoice` tidak boleh merupakan string kosong"}
	}

	if len(metodeBayar) == 0 {
		return Transaksi{}, oops.BadValues{
			Err: errors.New("`MetodeBayar` cannot be an empty string"),
			Msg: "`MetodeBayar` tidak boleh merupakan string kosong"}
	}

	t := Transaksi{
		ID:              transaksiId,
		IdUser:          idUser,
		IdAlamat:        idAlamat,
		HargaTotal:      uint(hargaTotal),
		KodeInvoice:     kodeInvoice,
		MetodeBayar:     metodeBayar,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		DetailTransaksi: detailTransaksi}
	return t, nil
}

func (t *Transaksi) WithAlamat(a Alamat) {
	t.Alamat = a
}
