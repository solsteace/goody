package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/catalog/internal/controller"
	"github.com/solsteace/goody/catalog/internal/lib/middleware"
)

func UseKategori(
	parent *fiber.Router,
	controller *controller.Kategori,
) {
	userContext := middleware.NewUserContext()
	adminChecker := middleware.NewAdminChecker()

	kategori := (*parent).Group("/produk")
	v1 := kategori.Group("/v1")
	v1.Get("/", controller.GetById)
	v1.Get("/:id", controller.GetById)

	v1.Use(userContext.Handle, adminChecker)
	v1.Post("/:id", controller.Create)
	v1.Put("/:id", controller.UpdateById)
	v1.Delete("/:id", controller.DeleteById)
}
