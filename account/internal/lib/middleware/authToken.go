package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/lib/token"
	"github.com/solsteace/goody/lib/token/payload"
)

type AuthToken struct {
	tokenHandler token.Handler[payload.Auth]
}

func NewAuthToken(handler token.Handler[payload.Auth]) AuthToken {
	return AuthToken{tokenHandler: handler}
}

func (a AuthToken) Handle(c *fiber.Ctx) error {
	token := c.Get("Authorization", "")
	if token == "" {
		return c.
			Status(http.StatusUnauthorized).
			JSON(fiber.Map{
				"message": "Token wasn't found in `Authorization` header",
			})
	}

	payload, err := a.tokenHandler.Decode(token)
	if err != nil {
		return err
	}

	c.Locals("Authorization", payload)
	return c.Next()
}
