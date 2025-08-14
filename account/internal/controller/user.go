package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/service"
	"github.com/solsteace/goody/lib/token/payload"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return UserController{service: service}
}

func (uc UserController) GetProfile(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*payload.AuthPayload)
	if !ok {
		return errors.New("Payload wasn't found on `Authorization` token")
	}

	result, err := uc.service.GetProfile(auth.UserId)
	if err != nil {
		return err
	}

	resPayload := map[string]any{
		"nama":          result.User.Nama,
		"no_telp":       result.User.NoTelp,
		"tanggal_lahir": result.User.TanggalLahir,
		"tentang":       result.User.Tentang,
		"pekerjaan":     result.User.Pekerjaan,
		"email":         result.User.Email,
		"id_provinsi":   <-result.Provinsi,
		"id_kota":       <-result.Kota}
	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    resPayload})
}

func (uc UserController) UpdateProfile(c *fiber.Ctx) error {
	reqPayload := new(struct {
		Nama         string `json:"nama"`
		TanggalLahir string `json:"tanggal_lahir"`
		Pekerjaan    string `json:"pekerjaan"`
		IdProvinsi   string `json:"id_provinsi"`
		IdKota       string `json:"id_kota"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	auth, ok := c.Locals("Authorization").(*payload.AuthPayload)
	if !ok {
		return errors.New("Authorization payload wasn't found")
	}

	tanggalLahir, err := time.Parse("02/01/2006", reqPayload.TanggalLahir)
	if err != nil {
		return err
	}

	result, err := uc.service.UpdateProfile(
		auth.UserId,
		reqPayload.Nama,
		tanggalLahir,
		reqPayload.Pekerjaan,
		reqPayload.IdProvinsi,
		reqPayload.IdKota)
	if err != nil {
		return err
	}

	resPayload := map[string]any{
		"nama":          result.User.Nama,
		"no_telp":       result.User.NoTelp,
		"tanggal_lahir": result.User.TanggalLahir,
		"tentang":       result.User.Tentang,
		"pekerjaan":     result.User.Pekerjaan,
		"email":         result.User.Email,
		"id_provinsi":   <-result.Provinsi,
		"id_kota":       <-result.Kota}
	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    resPayload})
}

func (uc UserController) ChangeCredentials(c *fiber.Ctx) error {
	reqPayload := new(struct {
		NoTelp        string `json:"no_telp"`
		Email         string `json:"email"`
		KataSandiLama string `json:"kata_sandi_lama"`
		KataSandiBaru string `json:"kata_sandi_baru"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	auth, ok := c.Locals("Authorization").(*payload.AuthPayload)
	if !ok {
		return errors.New("Payload wasn't found on `Authorization` token")
	}

	result, err := uc.service.ChangeCredentials(
		auth.UserId,
		reqPayload.NoTelp,
		reqPayload.Email,
		reqPayload.KataSandiLama,
		reqPayload.KataSandiBaru)
	if err != nil {
		return err
	}

	resPayload := map[string]any{
		"nama":          result.User.Nama,
		"no_telp":       result.User.NoTelp,
		"tanggal_lahir": result.User.TanggalLahir,
		"tentang":       result.User.Tentang,
		"pekerjaan":     result.User.Pekerjaan,
		"email":         result.User.Email,
		"id_provinsi":   <-result.Provinsi,
		"id_kota":       <-result.Kota}
	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    resPayload})
}
