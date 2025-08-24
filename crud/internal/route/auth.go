package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/controller"
	"github.com/solsteace/goody/crud/internal/util/middleware"
)

func UseAuth(
	parent *fiber.Router,
	controller *controller.Auth,
	authToken middleware.AuthChecker,
) {
	auth := (*parent).Group("/auth")
	v1 := auth.Group("/v1")

	v1.Post("/login", controller.Login)
	v1.Post("/register", controller.Register)
	v1.Use(authToken.Handle)
	v1.Get("/infer", controller.Infer)
}
