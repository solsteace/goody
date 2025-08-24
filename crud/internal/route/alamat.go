package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/controller"
)

func UseAlamat(
	parent *fiber.Router,
	controller *controller.Alamat,
) {
	alamat := (*parent).Group("/alamat")
	v1 := alamat.Group("/v1")

	v1.Get("/", controller.GetSelf)
	v1.Get("/:id", controller.GetById)
	v1.Post("/", controller.CreateForSelf)
	v1.Put("/:id", controller.UpdateById)
	v1.Delete("/:id", controller.DeleteById)
}
