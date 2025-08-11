package internal

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/lib/token/payload"
)

type updateProfilePayload struct {
	Nama         string `json:"nama"`
	KataSandi    string `json:"kata_sandi"`
	NoTelp       string `json:"no_telp"`
	TanggalLahir string `json:"tanggal_lahir"`
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `json:"email"`
	IdProvinsi   string `json:"id_provinsi"`
	IdKota       string `json:"id_kota"`
}

type UserController struct {
	service UserService
}

func NewUserController(service UserService) UserController {
	return UserController{service: service}
}

func (uc UserController) GetProfile(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*payload.AuthPayload)
	if !ok {
		return errors.New("Payload wasn't found on `Authorization` token")
	}

	userId := auth.UserId
	resPayload, err := uc.service.GetProfile(userId)
	if err != nil {
		return err
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    resPayload,
		})
}

func (uc UserController) UpdateProfile(c *fiber.Ctx) error {
	reqPayload := new(updateProfilePayload)
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

	userId := auth.UserId
	resPayload, err := uc.service.UpdateProfile(
		userId,
		reqPayload.Nama,
		reqPayload.KataSandi,
		reqPayload.NoTelp,
		tanggalLahir,
		reqPayload.Pekerjaan,
		reqPayload.Email,
		reqPayload.IdProvinsi,
		reqPayload.IdKota,
	)
	if err != nil {
		return err
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    resPayload,
		})
}
