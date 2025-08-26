package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/controller"
	"github.com/solsteace/goody/crud/internal/util/middleware"
)

func UseKategori(
	parent *fiber.Router,
	controller *controller.Kategori,
	authChecker middleware.AuthChecker,
) {
	adminChecker := middleware.NewAdminChecker()

	kategori := (*parent).Group("/kategori")
	kategori.Get("/", controller.GetMany)
	kategori.Get("/:id", controller.GetById)

	kategori.Use(authChecker.Handle, adminChecker.Handle)
	kategori.Post("/", controller.Create)
	kategori.Put("/:id", controller.UpdateById)
	kategori.Delete("/:id", controller.DeleteById)
}
