package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/controller"
	"github.com/solsteace/goody/crud/internal/util/middleware"
)

func UseTransaksi(
	parent *fiber.Router,
	controller *controller.Transaksi,
	authChecker middleware.AuthChecker,
) {
	transaksi := (*parent).Group("/trx")

	transaksi.Use(authChecker.Handle)
	transaksi.Get("/:id", controller.GetById)
	transaksi.Get("/", controller.GetSelf)
	transaksi.Post("/", controller.Create)
}
