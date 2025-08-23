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
	"github.com/solsteace/goody/lib/token"
)

type User struct {
	viewer  view.User
	service *service.User
}

func NewUser(service *service.User, viewer view.User) User {
	return User{viewer: viewer, service: service}
}

func (uc User) GetProfile(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		err := oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	result, err := uc.service.GetProfile(auth.UserId)
	if err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": fmt.Sprintf("Succeed to %s data", c.Method()),
			"errors":  nil,
			"data": fiber.Map{
				"user": uc.viewer.User(result.User),
			}})
}

func (uc User) UpdateProfile(c *fiber.Ctx) error {
	reqPayload := new(struct {
		Nama         string `json:"nama"`
		TanggalLahir string `json:"tanggal_lahir"`
		Pekerjaan    string `json:"pekerjaan"`
		IdProvinsi   string `json:"id_provinsi"`
		IdKota       string `json:"id_kota"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		err := oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	tanggalLahir, err := time.Parse("02/01/2006", reqPayload.TanggalLahir)
	if err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	result, err := uc.service.UpdateProfile(
		auth.UserId,
		reqPayload.Nama,
		tanggalLahir,
		reqPayload.Pekerjaan,
		reqPayload.IdProvinsi,
		reqPayload.IdKota)
	if err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data": fiber.Map{
				"user": uc.viewer.User(result.User),
			}})
}

func (uc User) ChangeCredentials(c *fiber.Ctx) error {
	reqPayload := new(struct {
		NoTelp        string `json:"no_telp"`
		Email         string `json:"email"`
		KataSandiLama string `json:"kata_sandi_lama"`
		KataSandiBaru string `json:"kata_sandi_baru"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	auth, ok := c.Locals("Authorization").(*token.Auth)
	if !ok {
		err := oops.Unauthorized{
			Err: errors.New("Payload wasn't found on `Authorization` token"),
			Msg: "Tidak ditemukan payload yang sesuai pada token"}
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	result, err := uc.service.ChangeCredentials(
		auth.UserId,
		reqPayload.NoTelp,
		reqPayload.Email,
		reqPayload.KataSandiLama,
		reqPayload.KataSandiBaru)
	if err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data": fiber.Map{
				"user": uc.viewer.User(result.User),
			}})
}
