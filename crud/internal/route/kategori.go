package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/controller"
	"github.com/solsteace/goody/crud/internal/util/middleware"
)

func UseKategori(
	parent *fiber.Router,
	controller *controller.Kategori,
) {
	adminChecker := middleware.NewAdminChecker()

	kategori := (*parent).Group("/produk")
	kategori.Get("/", controller.GetById)
	kategori.Get("/:id", controller.GetById)

	kategori.Use(adminChecker.Handle)
	kategori.Post("/:id", controller.Create)
	kategori.Put("/:id", controller.UpdateById)
	kategori.Delete("/:id", controller.DeleteById)
}
