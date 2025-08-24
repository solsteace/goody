package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/catalog/internal/controller"
	"github.com/solsteace/goody/catalog/internal/lib/middleware"
)

func UseToko(
	parent *fiber.Router,
	controller *controller.Toko,
) {
	userContext := middleware.NewUserContext()

	toko := (*parent).Group("/toko")
	v1 := toko.Group("/v1")

	v1.Get("/my", userContext.Handle, controller.GetSelf)
	v1.Get("/:id", controller.GetById)
	v1.Get("/", controller.GetMany)
	v1.Put("/:id", userContext.Handle, controller.UpdateById)
}
