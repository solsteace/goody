package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/lib/view"
	"github.com/solsteace/goody/account/internal/service"
	"github.com/solsteace/goody/lib/oops"
	"github.com/solsteace/goody/lib/oops/adapter"
	"github.com/solsteace/goody/lib/payload"
	"github.com/solsteace/goody/lib/token"
)

type Auth struct {
	service   *service.Auth
	authView  view.Auth
	payloader payload.Loader
}

func NewAuth(
	service *service.Auth,
	authView view.Auth,
	payloader payload.Loader,
) Auth {
	return Auth{service, authView, payloader}
}

func (ac Auth) Login(c *fiber.Ctx) error {
	reqPayload := new(struct {
		KataSandi string `json:"kata_sandi"`
		NoTelp    string `json:"no_telp"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(ac.payloader.Err(c.Method(), []error{err}))
	}

	result, err := ac.service.Login(reqPayload.NoTelp, reqPayload.KataSandi)
	if err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(ac.payloader.Err(c.Method(), []error{err}))
	}

	resPayload := ac.authView.Login(result.User, result.AccessToken, result.RefreshToken)
	return c.
		Status(http.StatusOK).
		JSON(ac.payloader.Ok(c.Method(), resPayload))
}

func (ac Auth) Register(c *fiber.Ctx) error {
	reqPayload := new(struct {
		Nama         string `json:"nama"`
		KataSandi    string `json:"kata_sandi"`
		NoTelp       string `json:"no_telp"`
		TanggalLahir string `json:"tanggal_lahir"`
		JenisKelamin string `json:"jenis_kelamin"`
		Tentang      string `json:"tentang"`
		Pekerjaan    string `json:"pekerjaan"`
		Email        string `json:"email"`
		IsAdmin      bool   `json:"is_admin"`
		IdProvinsi   string `json:"id_provinsi"`
		IdKota       string `json:"id_kota"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(ac.payloader.Err(c.Method(), []error{err}))
	}

	tanggalLahir, err := time.Parse("02/01/2006", reqPayload.TanggalLahir)
	if err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(ac.payloader.Err(c.Method(), []error{err}))
	}

	err = ac.service.Register(
		reqPayload.Nama,
		reqPayload.KataSandi,
		reqPayload.NoTelp,
		tanggalLahir,
		reqPayload.JenisKelamin,
		reqPayload.Tentang,
		reqPayload.Pekerjaan,
		reqPayload.Email,
		reqPayload.IsAdmin,
		reqPayload.IdProvinsi,
		reqPayload.IdKota)
	if err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(ac.payloader.Err(c.Method(), []error{err}))
	}

	resPayload := "Register Succeed"
	return c.
		Status(http.StatusCreated).
		JSON(ac.payloader.Ok(c.Method(), resPayload))
}

func (ac Auth) Infer(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		err := oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(ac.payloader.Err(c.Method(), []error{err}))
	}

	isAdmin := "0"
	if auth.IsAdmin {
		isAdmin = "1"
	}

	c.Response().Header.Add("X-User-Id", fmt.Sprintf("%d", auth.UserId))
	c.Response().Header.Add("X-User-IsAdmin", isAdmin)
	return c.SendStatus(http.StatusOK)
}
