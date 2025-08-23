package middleware

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/lib/oops"
	"github.com/solsteace/goody/lib/oops/adapter"
	"github.com/solsteace/goody/lib/token"
)

type AuthToken struct {
	tokenHandler token.Handler[token.Auth]
}

func NewAuthToken(handler token.Handler[token.Auth]) AuthToken {
	return AuthToken{tokenHandler: handler}
}

func (a AuthToken) Handle(c *fiber.Ctx) error {
	token := c.Get("Authorization", "")
	if token == "" {
		err := oops.Unauthorized{
			Err: errors.New("Token wasn't found in `Authorization` header"),
			Msg: "Token tidak ditemukan pada header `Authorization`"}
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	payload, err := a.tokenHandler.Decode(token)
	if err != nil {
		return c.
			Status(adapter.HttpStatusCode(err)).
			JSON(fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("Failed to %s data", c.Method()),
				"errors":  []string{err.Error()},
				"data":    ""})
	}

	c.Locals("Authorization", payload)
	return c.Next()
}
