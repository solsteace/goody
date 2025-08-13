package controller

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/service"
)

type loginPayload struct {
	KataSandi string `json:"kata_sandi"`
	NoTelp    string `json:"no_telp"`
}

type registerPayload struct {
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
}

type AuthController struct {
	service service.AuthService
}

func NewAuthController(as service.AuthService) AuthController {
	return AuthController{
		service: as,
	}
}

func (ac AuthController) Login(c *fiber.Ctx) error {
	reqPayload := new(loginPayload)
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	resData, err := ac.service.Login(reqPayload.NoTelp, reqPayload.KataSandi)
	if err != nil {
		return c.SendString(err.Error())
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    resData,
		})
}

func (ac AuthController) Register(c *fiber.Ctx) error {
	reqPayload := new(registerPayload)
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
		SendString("Register Succeed")
}
