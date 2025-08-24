package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/controller"
)

func UseAlamat(
	parent *fiber.Router,
	controller *controller.Alamat,
) {
	alamat := (*parent).Group("/alamat")
	alamatV1 := alamat.Group("/v1")
	alamatV1.Get("/", controller.GetSelf)
	alamatV1.Get("/:id", controller.GetById)
	alamatV1.Post("/", controller.CreateForSelf)
	alamatV1.Put("/:id", controller.UpdateById)
	alamatV1.Delete("/:id", controller.DeleteById)
}
