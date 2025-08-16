package route

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/controller"
	"github.com/solsteace/goody/lib/token"
	"github.com/solsteace/goody/lib/token/payload"
)

func UseUser(
	parent *fiber.Router,
	userController *controller.User,
	alamatController *controller.Alamat,
	tokenHandler token.Handler[payload.AuthPayload],
) {
	user := (*parent).Group("/user")
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

	UseAlamat(&user, alamatController)

	userV1 := user.Group("/v1")
	userV1.Get("/", userController.GetProfile)
	userV1.Put("/", userController.UpdateProfile)
	userV1.Patch("/credentials", userController.ChangeCredentials)
}
