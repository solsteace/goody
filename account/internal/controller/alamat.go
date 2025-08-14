package controller

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/service"
	"github.com/solsteace/goody/lib/token/payload"
)

type AlamatController struct {
	service service.AlamatService
}

func NewAlamatController(service service.AlamatService) AlamatController {
	return AlamatController{service: service}
}

func (ac AlamatController) GetSelf(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*payload.AuthPayload)
	if !ok {
		return errors.New("Payload wasn't found on `Authorization` token")
	}

	offset := uint(c.QueryInt("offset", 0))
	limit := uint(c.QueryInt("limit", 10))
	result, err := ac.service.GetSelf(auth.UserId, offset, limit)
	if err != nil {
		return err
	}

	resPayload := []struct {
		Id           uint   `json:"id"`
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp" `
		DetailAlamat string `json:"detail_alamat"`
	}{}
	for _, a := range result.Alamat {
		p := struct {
			Id           uint   `json:"id"`
			JudulAlamat  string `json:"judul_alamat"`
			NamaPenerima string `json:"nama_penerima"`
			NoTelp       string `json:"no_telp" `
			DetailAlamat string `json:"detail_alamat"`
		}{
			Id:           a.ID,
			JudulAlamat:  a.JudulAlamat,
			NamaPenerima: a.NamaPenerima,
			NoTelp:       a.NoTelp,
			DetailAlamat: a.DetailAlamat,
		}
		resPayload = append(resPayload, p)
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to GET data",
			"errors":  nil,
			"data":    resPayload,
		})
}

func (ac AlamatController) GetById(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*payload.AuthPayload)
	if !ok {
		return errors.New("Payload wasn't found on `Authorization` token")
	}

	alamatId, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	result, err := ac.service.GetById(auth.UserId, uint(alamatId))
	if err != nil {
		return err
	}

	resPayload := struct {
		Id           uint   `json:"id"`
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp" `
		DetailAlamat string `json:"detail_alamat"`
	}{
		Id:           result.Alamat.ID,
		JudulAlamat:  result.Alamat.JudulAlamat,
		NamaPenerima: result.Alamat.NamaPenerima,
		NoTelp:       result.Alamat.NoTelp,
		DetailAlamat: result.Alamat.DetailAlamat,
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

func (ac AlamatController) CreateForSelf(c *fiber.Ctx) error {
	reqPayload := new(struct {
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		DetailAlamat string `json:"detail_alamat"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	auth, ok := c.Locals("Authorization").(*payload.AuthPayload)
	if !ok {
		return errors.New("Payload wasn't found on `Authorization` token")
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

	resPayload := result.Alamat.ID
	return c.
		Status(http.StatusCreated).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    resPayload,
		})
}

func (ac AlamatController) UpdateById(c *fiber.Ctx) error {
	reqPayload := new(struct {
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		DetailAlamat string `json:"detail_alamat"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	auth, ok := c.Locals("Authorization").(*payload.AuthPayload)
	if !ok {
		return errors.New("Payload wasn't found on `Authorization` token")
	}

	alamatId, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	err = ac.service.UpdateById(
		auth.UserId,
		uint(alamatId),
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
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    resPayload,
		})
}

func (ac AlamatController) DeleteById(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(*payload.AuthPayload)
	if !ok {
		return errors.New("Payload wasn't found on `Authorization` token")
	}

	alamatId, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	if err = ac.service.DeleteById(auth.UserId, uint(alamatId)); err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}
