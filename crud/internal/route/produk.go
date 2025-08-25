package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/controller"
	"github.com/solsteace/goody/crud/internal/util/middleware"
)

func UseProduk(
	parent *fiber.Router,
	controller *controller.Produk,
	authChecker middleware.AuthChecker,
) {
	produk := (*parent).Group("/produk")
	v1 := produk.Group("/v1")

	v1.Get("/", controller.GetMany)
	v1.Get("/:id", controller.GetById)

	v1.Use(authChecker.Handle)
	v1.Post("/", controller.Create)
	v1.Put("/:id", controller.UpdateById)
	v1.Delete("/:id", controller.DeleteById)
}
