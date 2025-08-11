package internal

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/lib/token"
	"github.com/solsteace/goody/account/internal/lib/token/payload"
)

func registerUserRoutes(
	parent *fiber.Router,
	controller *UserController,
	tokenHandler token.Handler[payload.AuthPayload],
) {
	user := (*parent).Group("/user")

	// TODO: Move to shared lib
	user.Use(func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return errors.New("Token not found")
		}

		payload, err := tokenHandler.Decode(token)
		if err != nil {
			return err
		}

		c.Locals("Authorization", payload)
		return c.Next()
	})

	userV1 := user.Group("/v1")
	userV1.Get("/", controller.GetProfile)
	userV1.Put("/", controller.UpdateProfile)
}
