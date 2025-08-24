package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/controller"
	"github.com/solsteace/goody/account/internal/lib/middleware"
)

func UseAuth(
	parent *fiber.Router,
	controller *controller.Auth,
	authToken middleware.AuthToken,
) {
	auth := (*parent).Group("/auth")
	authV1 := auth.Group("/v1")
	authV1.Post("/login", controller.Login)
	authV1.Post("/register", controller.Register)

	authV1.Use(authToken.Handle)
	authV1.Get("/infer", controller.Infer)
}
