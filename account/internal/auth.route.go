package internal

import (
	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(parent *fiber.Router, controller *AuthController) {
	auth := (*parent).Group("/auth")
	authV1 := auth.Group("/v1")
	authV1.Post("/login", controller.login)
	authV1.Post("/register", controller.register)
}
