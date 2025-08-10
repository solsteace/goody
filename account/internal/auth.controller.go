package internal

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
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
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `json:"email"`
	IdProvinsi   string `json:"id_provinsi"`
	IdKota       string `json:"id_kota"`
}

type AuthController struct {
	service AuthService
}

func NewAuthController(as AuthService) AuthController {
	return AuthController{
		service: as,
	}
}

func (ac AuthController) login(c *fiber.Ctx) error {
	reqPayload := new(loginPayload)
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	resData, err := ac.service.login(reqPayload.NoTelp, reqPayload.KataSandi)
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

func (ac AuthController) register(c *fiber.Ctx) error {
	reqPayload := new(registerPayload)
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	birthdate, err := time.Parse("02/01/2006", reqPayload.TanggalLahir)
	if err != nil {
		return err
	}

	err = ac.service.register(
		reqPayload.Nama,
		reqPayload.KataSandi,
		reqPayload.NoTelp,
		birthdate,
		reqPayload.Pekerjaan,
		reqPayload.Email,
		reqPayload.IdProvinsi,
		reqPayload.IdKota)
	if err != nil {
		return err
	}

	return c.
		Status(http.StatusCreated).
		SendString("Register Succeed")
}
