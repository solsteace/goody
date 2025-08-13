package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/controller"
)

func registerAlamatRoutes(
	parent *fiber.Router,
	alamatController *controller.AlamatController,
) {
	alamat := (*parent).Group("/alamat")
	alamatV1 := alamat.Group("/v1")
	alamatV1.Get("/", alamatController.GetSelf)
	alamatV1.Get("/:id", alamatController.GetById)
	alamatV1.Post("/", alamatController.CreateForSelf)
	alamatV1.Put("/:id", alamatController.UpdateById)
	alamatV1.Delete("/:id", alamatController.DeleteById)
}
