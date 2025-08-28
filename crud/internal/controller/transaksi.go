package controller

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/repository"
	"github.com/solsteace/goody/crud/internal/service"
	"github.com/solsteace/goody/crud/internal/util/view"
	"github.com/solsteace/goody/lib/oops"
	"github.com/solsteace/goody/lib/payload"
	"github.com/solsteace/goody/lib/token"
)

type Transaksi struct {
	service   service.Transaksi
	viewer    view.Transaksi
	payloader payload.Loader
}

func NewTransaksi(
	service service.Transaksi,
	viewer view.Transaksi,
	payloader payload.Loader,
) Transaksi {
	return Transaksi{service, viewer, payloader}
}

func (tc Transaksi) GetSelf(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}
	page := c.QueryInt("page")
	limit := c.QueryInt("limit")

	result, err := tc.service.GetMany(&auth.UserId, page, limit)
	if err != nil {
		return err
	}

	resPayload := tc.viewer.ManyTransaksi(result.Transaksi)
	return c.
		Status(http.StatusOK).
		JSON(tc.payloader.Ok(c.Method(), resPayload))
}

func (tc Transaksi) GetById(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}

	idTransaksi, _ := c.ParamsInt("id", 0)
	result, err := tc.service.GetById(auth.UserId, uint(idTransaksi))
	if err != nil {
		return err
	}

	resPayload := tc.viewer.Transaksi(result.Transaksi)
	return c.
		Status(http.StatusOK).
		JSON(tc.payloader.Ok(c.Method(), resPayload))
}

func (tc Transaksi) Create(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}

	reqPayload := new(struct {
		IdAlamat        uint   `json:"alamat_kirim"`
		MetodeBayar     string `json:"method_bayar"`
		DetailTransaksi []struct {
			IdProduk  uint `json:"product_id"`
			Kuantitas int  `json:"kuantitas"`
		} `json:"detail_trx"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	detailTransaksi := []repository.DetailTransaksiEntry{}
	for _, dt := range reqPayload.DetailTransaksi {
		entry, err := repository.NewDetailTransaksiEntry(dt.IdProduk, dt.Kuantitas)
		if err != nil {
			return err
		}
		detailTransaksi = append(detailTransaksi, entry)
	}
	result, err := tc.service.Create(
		auth.UserId,
		reqPayload.IdAlamat,
		reqPayload.MetodeBayar,
		detailTransaksi)
	if err != nil {
		return err
	}

	resPayload := tc.viewer.Transaksi(result.Transaksi)
	return c.
		Status(http.StatusCreated).
		JSON(tc.payloader.Ok(c.Method(), resPayload))
}
