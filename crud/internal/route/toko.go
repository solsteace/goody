package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/controller"
	"github.com/solsteace/goody/crud/internal/util/middleware"
)

func UseToko(
	parent *fiber.Router,
	controller *controller.Toko,
	authChecker middleware.AuthChecker,
) {
	toko := (*parent).Group("/toko")
	toko.Get("/my", controller.GetSelf)
	toko.Get("/:id", controller.GetById)
	toko.Get("/", controller.GetMany)

	toko.Use(authChecker.Handle)
	toko.Put("/:id", controller.UpdateById)
}
