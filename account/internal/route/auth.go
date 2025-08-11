package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/controller"
)

func RegisterAuthRoutes(parent *fiber.Router, controller *controller.AuthController) {
	auth := (*parent).Group("/auth")
	authV1 := auth.Group("/v1")
	authV1.Post("/login", controller.Login)
	authV1.Post("/register", controller.Register)
}
