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
	produk.Get("/", controller.GetMany)
	produk.Get("/:id", controller.GetById)

	produk.Use(authChecker.Handle)
	produk.Post("/", controller.Create)
	produk.Put("/:id", controller.UpdateById)
	produk.Delete("/:id", controller.DeleteById)
}
