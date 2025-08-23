package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/lib/view"
	"github.com/solsteace/goody/account/internal/service"
	"github.com/solsteace/goody/lib/oops"
	"github.com/solsteace/goody/lib/oops/adapter"
	"github.com/solsteace/goody/lib/token"
)

type Alamat struct {
	service *service.Alamat
	viewer  view.Alamat
}

func NewAlamat(
	service *service.Alamat,
	viewer view.Alamat,
) Alamat {
	return Alamat{service: service, viewer: viewer}
}

func (ac Alamat) GetSelf(c *fiber.Ctx) error {
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

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	judul := c.Query("judul_alamat", "")
	result, err := ac.service.GetSelf(auth.UserId, judul, page, limit)
	if err != nil {
		return err
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to GET data",
			"errors":  nil,
			"data":    ac.viewer.ManyAlamat(result.Alamat)})
}

func (ac Alamat) GetById(c *fiber.Ctx) error {
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

	alamatId, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	result, err := ac.service.GetById(auth.UserId, uint(alamatId))
	if err != nil {
		return err
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to POST data",
			"errors":  nil,
			"data":    ac.viewer.Alamat(result.Alamat)})
}

func (ac Alamat) CreateForSelf(c *fiber.Ctx) error {
	reqPayload := new(struct {
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		DetailAlamat string `json:"detail_alamat"`
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
			"data":    resPayload})
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

	alamatId, err := c.ParamsInt("id")
	if err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	err = ac.service.UpdateById(
		auth.UserId,
		uint(alamatId),
		reqPayload.JudulAlamat,
		reqPayload.NamaPenerima,
		reqPayload.NoTelp,
		reqPayload.DetailAlamat)
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
			"data":    ""})
}

func (ac Alamat) DeleteById(c *fiber.Ctx) error {
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

	alamatId, err := c.ParamsInt("id")
	if err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	if err = ac.service.DeleteById(auth.UserId, uint(alamatId)); err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	return c.
		Status(http.StatusCreated).
		JSON(fiber.Map{
			"status":  true,
			"message": fmt.Sprintf("Succeed to %s data", c.Method()),
			"errors":  nil,
			"data":    ""})
}
