package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/lib/oops"
	"github.com/solsteace/goody/lib/token"
)

type AuthChecker struct {
	tokenHandler token.Handler[token.Auth]
}

func NewAuthToken(handler token.Handler[token.Auth]) AuthChecker {
	return AuthChecker{tokenHandler: handler}
}

func (a AuthChecker) Handle(c *fiber.Ctx) error {
	token := c.Get("Authorization", "")
	if token == "" {
		return oops.Unauthorized{
			Err: errors.New("Token wasn't found in `Authorization` header"),
			Msg: "Token tidak ditemukan pada header `Authorization`"}
	}

	payload, err := a.tokenHandler.Decode(token)
	if err != nil {
		return err
	}

	c.Locals("Authorization", payload)
	return c.Next()
}
