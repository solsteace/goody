package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/controller"
)

func UseToko(
	parent *fiber.Router,
	controller *controller.Toko,
) {
	toko := (*parent).Group("/toko")
	v1 := toko.Group("/v1")

	v1.Get("/my", controller.GetSelf)
	v1.Get("/:id", controller.GetById)
	v1.Get("/", controller.GetMany)
	v1.Put("/:id", controller.UpdateById)
}
