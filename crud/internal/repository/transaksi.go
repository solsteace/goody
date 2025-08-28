package repository

import (
	"errors"

	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/lib/oops"
)

const transaksiDefaultPageSize = 10

type Transaksi interface {
	GetMany(q transaksiQueryParams) ([]domain.Transaksi, error)
	GetById(id uint) (domain.Transaksi, error)

	// Special case: Since the creation of `transaksi` depends on the state of
	// other tables in order to be created (in this case, it's product),
	// we pass "the materials" to build them instead of the instance itself.
	Create(
		idUser,
		idToko,
		idAlamat uint,
		kodeInvoice,
		metodeBayar string,
		detailTransaksiEntry []DetailTransaksiEntry,
	) (uint, error)
}

type DetailTransaksiEntry struct {
	idProduk  uint
	kuantitas uint
}

func NewDetailTransaksiEntry(idProduk uint, kuantitas int) (DetailTransaksiEntry, error) {
	if kuantitas < 0 {
		return DetailTransaksiEntry{}, oops.BadValues{
			Err: errors.New("`Kuantitas` must be a positive number"),
			Msg: "Kuantitas harus bernilai positif"}
	}

	return DetailTransaksiEntry{idProduk, uint(kuantitas)}, nil
}

type transaksiQueryParams struct {
	page   int
	limit  int
	idUser uint
}

func (param transaksiQueryParams) offset() int {
	if param.page < 1 {
		return 0
	}
	return (param.page - 1) * param.limit
}

func NewTransaksiQueryParams(
	page,
	limit int,
	idUser *uint,
) transaksiQueryParams {
	actualPage := 0
	if page != 0 {
		actualPage = page
	}

	actualLimit := transaksiDefaultPageSize
	if limit != 0 {
		actualLimit = limit
	}

	var actualIdUser uint = 0
	if idUser != nil {
		actualIdUser = *idUser
	}

	return transaksiQueryParams{
		page:   actualPage,
		limit:  actualLimit,
		idUser: actualIdUser}
}
