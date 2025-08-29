package controller

import (
	"errors"
	"mime/multipart"
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
		page, limit, nama, maxHarga, minHarga, categoryId, tokoId)
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
	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}

	reqPayload := new(struct {
		NamaProduk    string `form:"nama_produk"`
		HargaReseller int    `form:"harga_reseller"`
		HargaKonsumer int    `form:"harga_konsumen"`
		Stok          int    `form:"stok"`
		Deskripsi     string `form:"deskripsi"`
		IdKategori    uint   `form:"category_id"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	foto, ok := form.File["photos"]
	if !ok {
		foto = []*multipart.FileHeader{}
	}

	for _, f := range foto {
		switch f.Header["Content-Type"][0] {
		case "image/jpeg", "image/webp", "image/png": // Let's say we only accept jpg, jpeg, webp, and png
		default:
			return oops.BadRequest{
				Err: errors.New("Mime type should be either image/jpg, image/webp, or image/png"),
				Msg: "File yang dikirim harus dalam format .jpg, .jpeg, .webp, atau .png"}
		}
	}

	result, err := pc.service.Create(
		auth.UserId,
		reqPayload.IdKategori,
		reqPayload.NamaProduk,
		reqPayload.HargaReseller,
		reqPayload.HargaKonsumer,
		reqPayload.Stok,
		reqPayload.Deskripsi,
		foto,
		c.SaveFile) // The storing mechanism is coupled to the framework :/
	if err != nil {
		return err
	}

	resPayload := pc.viewer.Produk(result.Produk)
	return c.
		Status(http.StatusCreated).
		JSON(pc.payloader.Ok(c.Method(), resPayload))
}

func (pc Produk) UpdateById(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}
	reqPayload := new(struct {
		NamaProduk    string `form:"nama_produk"`
		HargaReseller int    `form:"harga_reseller"`
		HargaKonsumer int    `form:"harga_konsumen"`
		Stok          int    `form:"stok"`
		Deskripsi     string `form:"deskripsi"`
		IdKategori    uint   `form:"category_id"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	foto, ok := form.File["photos"]
	if !ok {
		foto = []*multipart.FileHeader{}
	}

	for _, f := range foto {
		switch f.Header["Content-Type"][0] {
		case "image/jpeg", "image/webp", "image/png": // Let's say we only accept jpg, jpeg, webp, and png
		default:
			return oops.BadRequest{
				Err: errors.New("Mime type should be either image/jpg, image/webp, or image/png"),
				Msg: "File yang dikirim harus dalam format .jpg, .jpeg, .webp, atau .png"}
		}
	}

	idProduk, _ := c.ParamsInt("id", 0)
	err = pc.service.UpdateById(
		auth.UserId,
		uint(idProduk),
		reqPayload.IdKategori,
		reqPayload.NamaProduk,
		reqPayload.HargaReseller,
		reqPayload.HargaKonsumer,
		reqPayload.Stok,
		reqPayload.Deskripsi,
		foto,
		c.SaveFile) // The storing mechanism is coupled to the framework :/
	if err != nil {
		return err
	}

	resPayload := ""
	return c.
		Status(http.StatusOK).
		JSON(pc.payloader.Ok(c.Method(), resPayload))
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
