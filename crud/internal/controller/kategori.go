package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/service"
	"github.com/solsteace/goody/crud/internal/util/view"
	"github.com/solsteace/goody/lib/payload"
)

type Kategori struct {
	service   service.Kategori
	viewer    view.Kategori
	payloader payload.Loader
}

func NewKategori(
	service service.Kategori,
	viewer view.Kategori,
	payloader payload.Loader,
) Kategori {
	return Kategori{service, viewer, payloader}
}

func (kc Kategori) GetMany(c *fiber.Ctx) error {
	page := c.QueryInt("page", 0)
	limit := c.QueryInt("limit", 10)

	result, err := kc.service.GetMany(page, limit)
	if err != nil {
		return err
	}

	resPayload := kc.viewer.ManyKategori(result.Kategori)
	return c.
		Status(http.StatusOK).
		JSON(kc.payloader.Ok(c.Method(), resPayload))
}

func (kc Kategori) GetById(c *fiber.Ctx) error {
	idKategori, _ := c.ParamsInt("id")
	result, err := kc.service.GetById(uint(idKategori))
	if err != nil {
		return err
	}

	resPayload := kc.viewer.Kategori(result.Kategori)
	return c.
		Status(http.StatusOK).
		JSON(kc.payloader.Ok(c.Method(), resPayload))
}

func (kc Kategori) Create(c *fiber.Ctx) error {
	reqPayload := new(struct {
		Nama string `json:"nama_category"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	result, err := kc.service.Create(reqPayload.Nama)
	if err != nil {
		return err
	}

	resPayload := kc.viewer.Kategori(result.Kategori)
	return c.
		Status(http.StatusCreated).
		JSON(kc.payloader.Ok(c.Method(), resPayload))
}

func (kc Kategori) UpdateById(c *fiber.Ctx) error {
	reqPayload := new(struct {
		Nama string `json:"nama_category"`
	})
	if err := c.BodyParser(reqPayload); err != nil {
		return err
	}

	idKategori, _ := c.ParamsInt("id", 0)
	if err := kc.service.UpdateById(uint(idKategori), reqPayload.Nama); err != nil {
		return err
	}

	resPayload := ""
	return c.
		Status(http.StatusOK).
		JSON(kc.payloader.Ok(c.Method(), resPayload))
}

func (kc Kategori) DeleteById(c *fiber.Ctx) error {
	idKategori, _ := c.ParamsInt("id", 0)
	if err := kc.service.DeleteById(uint(idKategori)); err != nil {
		return err
	}

	resPayload := ""
	return c.
		Status(http.StatusOK).
		JSON(kc.payloader.Ok(c.Method(), resPayload))
}
