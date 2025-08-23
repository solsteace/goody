package controller

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/lib/view"
	"github.com/solsteace/goody/account/internal/service"
	"github.com/solsteace/goody/lib/oops"
	"github.com/solsteace/goody/lib/payload"
	"github.com/solsteace/goody/lib/token"
)

type Alamat struct {
	service    *service.Alamat
	alamatView view.Alamat
	payloader  payload.Loader
}

func NewAlamat(
	service *service.Alamat,
	alamatView view.Alamat,
	payloader payload.Loader,
) Alamat {
	return Alamat{service, alamatView, payloader}
}

func (ac Alamat) GetSelf(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	judul := c.Query("judul_alamat", "")
	result, err := ac.service.GetSelf(auth.UserId, judul, page, limit)
	if err != nil {
		return err
	}

	resPayload := ac.alamatView.ManyAlamat(result.Alamat)
	return c.
		Status(http.StatusOK).
		JSON(ac.payloader.Ok(c.Method(), resPayload))
}

func (ac Alamat) GetById(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}

	idAlamat, _ := c.ParamsInt("id", 0)
	result, err := ac.service.GetById(auth.UserId, uint(idAlamat))
	if err != nil {
		return err
	}

	resPayload := ac.alamatView.Alamat(result.Alamat)
	return c.
		Status(http.StatusOK).
		JSON(ac.payloader.Ok(c.Method(), resPayload))
}

func (ac Alamat) CreateForSelf(c *fiber.Ctx) error {
	reqPayload := new(struct {
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		DetailAlamat string `json:"detail_alamat"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}

	result, err := ac.service.CreateForSelf(
		auth.UserId,
		reqPayload.JudulAlamat,
		reqPayload.NamaPenerima,
		reqPayload.NoTelp,
		reqPayload.DetailAlamat)
	if err != nil {
		return err
	}

	resPayload := ac.alamatView.Alamat(result.Alamat)
	return c.
		Status(http.StatusCreated).
		JSON(ac.payloader.Ok(c.Method(), resPayload))
}

func (ac Alamat) UpdateById(c *fiber.Ctx) error {
	reqPayload := new(struct {
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		DetailAlamat string `json:"detail_alamat"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}

	idAlamat, err := c.ParamsInt("id", 0)
	err = ac.service.UpdateById(
		auth.UserId,
		uint(idAlamat),
		reqPayload.JudulAlamat,
		reqPayload.NamaPenerima,
		reqPayload.NoTelp,
		reqPayload.DetailAlamat)
	if err != nil {
		return err
	}

	resPayload := ""
	return c.
		Status(http.StatusOK).
		JSON(ac.payloader.Ok(c.Method(), resPayload))
}

func (ac Alamat) DeleteById(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		return oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
	}

	idAlamat, _ := c.ParamsInt("id", 0)
	if err := ac.service.DeleteById(auth.UserId, uint(idAlamat)); err != nil {
		return err
	}

	resPayload := ""
	return c.
		Status(http.StatusCreated).
		JSON(ac.payloader.Ok(c.Method(), resPayload))
}
