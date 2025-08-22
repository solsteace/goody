package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/lib/view"
	"github.com/solsteace/goody/account/internal/service"
	"github.com/solsteace/goody/lib/token/payload"
)

type Auth struct {
	service *service.Auth
	viewer  view.Auth
}

func NewAuth(service *service.Auth, viewer view.Auth) Auth {
	return Auth{viewer: viewer, service: service}
}

func (ac Auth) Login(c *fiber.Ctx) error {
	reqPayload := new(struct {
		KataSandi string `json:"kata_sandi"`
		NoTelp    string `json:"no_telp"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	result, err := ac.service.Login(reqPayload.NoTelp, reqPayload.KataSandi)
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data": ac.viewer.Login(
				result.User, result.AccessToken, result.RefreshToken),
		})
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
		return err
	}

	tanggalLahir, err := time.Parse("02/01/2006", reqPayload.TanggalLahir)
	if err != nil {
		return err
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
		return err
	}

	return c.
		Status(http.StatusCreated).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    "Register Succeed",
		})
}

func (ac Auth) Infer(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*payload.Auth)
	if !ok {
		return errors.New("Valid payload wasn't found on token")
	}

	isAdmin := "0"
	if auth.IsAdmin {
		isAdmin = "1"
	}

	c.Response().Header.Add("X-User-Id", fmt.Sprintf("%d", auth.UserId))
	c.Response().Header.Add("X-User-IsAdmin", isAdmin)
	return c.SendStatus(http.StatusOK)
}
