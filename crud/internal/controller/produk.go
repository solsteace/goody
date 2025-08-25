package controller

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/service"
	"github.com/solsteace/goody/crud/internal/util/view"
	"github.com/solsteace/goody/lib/oops"
	"github.com/solsteace/goody/lib/payload"
	"github.com/solsteace/goody/lib/token"
)

type Produk struct {
	service   service.Produk
	viewer    view.Produk
	payloader payload.Loader
}

func NewProduk(
	service service.Produk,
	viewer view.Produk,
	payloader payload.Loader,
) Produk {
	return Produk{service, viewer, payloader}
}

func (pc Produk) GetMany(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	limit := c.QueryInt("limit")
	nama := c.Query("nama_produk")
	minHarga := c.QueryInt("min_harga")
	maxHarga := c.QueryInt("max_harga")
	tokoId := c.QueryInt("toko_id")
	categoryId := c.QueryInt("category_id")

	result, err := pc.service.GetMany(
		page, limit, nama, minHarga, maxHarga, tokoId, categoryId)
	if err != nil {
		return err
	}

	resPayload := pc.viewer.ManyProduk(result.Produk)
	return c.
		Status(http.StatusOK).
		JSON(pc.payloader.Ok(c.Method(), resPayload))
}

func (pc Produk) GetById(c *fiber.Ctx) error {
	idProduk, _ := c.ParamsInt("id", 0)

	result, err := pc.service.GetById(uint(idProduk))
	if err != nil {
		return err
	}

	resPayload := pc.viewer.Produk(result.Produk)
	return c.
		Status(http.StatusOK).
		JSON(pc.payloader.Ok(c.Method(), resPayload))
}

func (pc Produk) Create(c *fiber.Ctx) error {
	reqPayload := new(struct {
		NamaProduk    string `json:"nama_produk"`
		HargaReseller int    `json:"harga_reseller"`
		HargaKonsumer int    `json:"harga_konsumen"`
		Stok          int    `json:"Stok"`
	})
	if err := c.BodyParser(&reqPayload); err != nil {
		return err
	}

	// Take photo

	return c.SendStatus(http.StatusNotImplemented)
}

func (pc Produk) UpdateById(c *fiber.Ctx) error {
	_, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}

	return c.SendStatus(http.StatusNotImplemented)
}

func (pc Produk) DeleteById(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}
	idProduk, _ := c.ParamsInt("id", 0)

	err := pc.service.DeleteById(auth.UserId, uint(idProduk))
	if err != nil {
		return err
	}

	return c.
		Status(http.StatusOK).
		JSON(pc.payloader.Ok(c.Method(), nil))
}
