package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/service"
	"github.com/solsteace/goody/lib/token/payload"
)

type Auth struct {
	service service.Auth
}

func NewAuth(as service.Auth) Auth {
	return Auth{
		service: as,
	}
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

	resData := map[string]any{
		"nama":          result.User.Nama,
		"no_telp":       result.User.NoTelp,
		"tanggal_lahir": result.User.TanggalLahir,
		"tentang":       result.User.Tentang,
		"pekerjaan":     result.User.Pekerjaan,
		"email":         result.User.Email,
		"id_provinsi":   <-result.Provinsi,
		"id_kota":       <-result.Kota,
		"token":         result.AccessToken}
	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    resData,
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
		SendString("Register Succeed")
}

func (ac Auth) Infer(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*payload.AuthPayload)
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
