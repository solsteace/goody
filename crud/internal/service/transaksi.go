package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/solsteace/goody/crud/internal/domain"
	"github.com/solsteace/goody/crud/internal/repository"
	"github.com/solsteace/goody/lib/oops"
)

type Transaksi struct {
	repo       repository.Transaksi
	alamatRepo repository.Alamat
	tokoRepo   repository.Toko
}

func NewTransaksi(
	repo repository.Transaksi,
	alamatRepo repository.Alamat,
	tokoRepo repository.Toko,
) Transaksi {
	return Transaksi{repo, alamatRepo, tokoRepo}
}

func (ts Transaksi) GetMany(idUser *uint, page, limit int) (
	*struct{ Transaksi []domain.Transaksi }, error,
) {
	result := new(struct{ Transaksi []domain.Transaksi })

	var actualIdUser *uint = nil
	if idUser != nil {
		actualIdUser = idUser
	}

	queryParams := repository.NewTransaksiQueryParams(page, limit, actualIdUser)
	transaksi, err := ts.repo.GetMany(queryParams)
	if err != nil {
		return result, err
	}

	result.Transaksi = transaksi
	return result, nil
}

func (ts Transaksi) GetById(idUser, id uint) (
	*struct{ Transaksi domain.Transaksi }, error,
) {
	result := new(struct{ Transaksi domain.Transaksi })

	transaksi, err := ts.repo.GetById(id)
	if err != nil {
		return result, err
	}

	if transaksi.IdUser != idUser {
		return result, oops.Forbidden{
			Err: errors.New(fmt.Sprintf("User(id:%d) doesn't own Transaksi(id:%d)",
				idUser, id)),
			Msg: "Anda tidak memiliki akses transaksi ini"}
	}

	result.Transaksi = transaksi
	return result, nil
}

func (ts Transaksi) Create(
	idUser uint,
	idAlamat uint,
	metodeBayar string,
	detailTransaksiEntry []repository.DetailTransaksiEntry,
) (*struct{ Transaksi domain.Transaksi }, error) {
	result := new(struct{ Transaksi domain.Transaksi })

	if len(metodeBayar) == 0 {
		return result, oops.BadValues{
			Err: errors.New("`MetodeBayar` cannot be an empty string"),
			Msg: "`MetodeBayar` tidak boleh merupakan string kosong"}
	}

	alamat, err := ts.alamatRepo.GetById(idAlamat)
	if err != nil {
		return result, err
	}
	if alamat.UserId != idUser {
		return result, oops.Forbidden{
			Err: errors.New(
				fmt.Sprintf("User(id:%d) doesn't own Alamat(id:%d)", idUser, idAlamat)),
			Msg: "Anda tidak dapat menggunakan alamat user lain untuk bertransaksi"}
	}

	toko, err := ts.tokoRepo.GetByOwnerId(idUser)
	if err != nil {
		return result, err
	}

	kodeInvoice := fmt.Sprintf("INV-%s", uuid.New().String())
	idTransaksi, err := ts.repo.Create(
		idUser,
		toko.ID,
		idAlamat,
		kodeInvoice,
		metodeBayar,
		detailTransaksiEntry)
	if err != nil {
		return result, err
	}

	transaksi, err := ts.repo.GetById(idTransaksi)
	if err != nil {
		return result, err
	}

	result.Transaksi = transaksi
	return result, nil
}
