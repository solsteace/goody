package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/catalog/internal/controller"
	"github.com/solsteace/goody/catalog/internal/lib/middleware"
)

func UseProduk(parent *fiber.Router, controller *controller.Produk) {
	produk := (*parent).Group("/produk")
	v1 := produk.Group("/v1")

	v1.Get("/", controller.GetMany)
	v1.Get("/:id", controller.GetById)

	v1.Use(middleware.UserContext)
	v1.Get("/my", controller.GetSelf)
	v1.Put("/", controller.UpdateById)
}
