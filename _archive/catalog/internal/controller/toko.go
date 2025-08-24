package controller

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/catalog/internal/lib/view"
	"github.com/solsteace/goody/catalog/internal/service"
	"github.com/solsteace/goody/lib/oops"
	"github.com/solsteace/goody/lib/payload"
)

type Toko struct {
	service   service.Toko
	viewer    view.Toko
	payloader payload.Loader
}

func NewToko(service service.Toko, viewer view.Toko) Toko {
	return Toko{viewer: viewer, service: service}
}

func (tc Toko) GetMany(c *fiber.Ctx) error {
	page := c.QueryInt("page", 0)
	limit := c.QueryInt("limit", 10)
	_ = c.Query("name", "")

	result, err := tc.service.GetMany(page, limit)
	if err != nil {
		return err
	}

	resPayload := tc.viewer.ManyToko(result.Toko)
	return c.
		Status(http.StatusOK).
		JSON(tc.payloader.Ok(c.Method(), resPayload))
}

func (tc Toko) GetSelf(c *fiber.Ctx) error {
	idUser, ok := c.Locals("userId").(uint)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Couldn't extract `userId`"),
			Msg: "Data `userId` tidak ditemukan "}
	}

	result, err := tc.service.GetByOwnerId(idUser)
	if err != nil {
		return err
	}

	resPayload := tc.viewer.Toko(result.Toko)
	return c.
		Status(http.StatusOK).
		JSON(tc.payloader.Ok(c.Method(), resPayload))
}

func (tc Toko) GetById(c *fiber.Ctx) error {
	idToko, _ := c.ParamsInt("id", 0)
	result, err := tc.service.GetById(uint(idToko))
	if err != nil {
		return err
	}

	resPayload := tc.viewer.Toko(result.Toko)
	return c.
		Status(http.StatusOK).
		JSON(tc.payloader.Ok(c.Method(), resPayload))
}

func (tc Toko) UpdateById(c *fiber.Ctx) error {
	idUser, ok := c.Locals("userId").(uint)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Couldn't extract `userId`"),
			Msg: "Data `userId` tidak ditemukan "}
	}

	reqPayload := new(struct {
		NamaToko string `json:"nama_toko"`
	})
	if err := c.BodyParser(&reqPayload); err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	formFiles, ok := form.File["photo"]
	if !ok || len(formFiles) == 0 {
		return oops.BadRequest{
			Err: errors.New("No file found in `photo`"),
			Msg: "Tidak ada file ditemukan pada `photo`"}
	}
	foto := formFiles[0]
	switch foto.Header["Content-Type"][0] {
	case "image/jpeg", "image/webp", "image/png": // Let's say we only accept jpg, jpeg, webp, and png
	default:
		return oops.BadRequest{
			Err: errors.New("Mime type should be either image/jpg, image/webp, or image/png"),
			Msg: "File yang dikirim harus dalam format .jpg, .jpeg, .webp, atau .png"}
	}

	// A lil' stinks, but this'll do. We'll refactor it later somehow
	idToko, _ := c.ParamsInt("id", 0)
	err = tc.service.UpdateById(
		uint(idToko),
		idUser,
		reqPayload.NamaToko,
		foto,
		func(savePath string) error {
			return c.SaveFile(foto, savePath)
		})
	if err != nil {
		return err
	}

	return c.
		Status(http.StatusOK).
		JSON(tc.payloader.Ok(c.Method(), ""))
}
