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
	toko.Get("/my", controller.GetSelf)
	toko.Get("/:id", controller.GetById)
	toko.Get("/", controller.GetMany)
	toko.Put("/:id", controller.UpdateById)
}
