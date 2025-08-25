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
	alamat.Get("/", controller.GetSelf)
	alamat.Get("/:id", controller.GetById)
	alamat.Post("/", controller.CreateForSelf)
	alamat.Put("/:id", controller.UpdateById)
	alamat.Delete("/:id", controller.DeleteById)
}
