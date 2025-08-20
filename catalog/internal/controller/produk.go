package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/catalog/internal/service"
)

type Produk struct {
	service service.Produk
}

func NewProduk(service service.Produk) Produk {
	return Produk{service: service}
}

func (pc Produk) GetSelf(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNotImplemented)
}

func (pc Produk) GetMany(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNotImplemented)
}

func (pc Produk) GetById(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNotImplemented)
}

func (pc Produk) Create(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNotImplemented)
}

func (pc Produk) UpdateById(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNotImplemented)
}

func (pc Produk) DeleteById(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusNotImplemented)
}
