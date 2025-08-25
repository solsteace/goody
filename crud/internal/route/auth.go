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
	auth.Post("/login", controller.Login)
	auth.Post("/register", controller.Register)
	auth.Use(authToken.Handle)
	auth.Get("/infer", controller.Infer)
}
