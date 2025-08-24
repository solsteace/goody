package middleware

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/lib/oops"
)

type userContext struct{}

func NewUserContext() userContext {
	return userContext{}
}

func (uc userContext) Handle(c *fiber.Ctx) error {
	userId, err := strconv.ParseUint(c.Get("X-User-Id"), 10, 64)
	if err != nil {
		return err
	}
	c.Locals("userId", uint(userId))

	if isAdmin := c.Get("X-User-IsAdmin"); isAdmin == "0" {
		c.Locals("isAdmin", false)
	} else if isAdmin == "1" {
		c.Locals("isAdmin", true)
	} else {
		return oops.BadRequest{
			Err: errors.New("`X-User-IsAdmin` value should be either '1' or '2'"),
			Msg: "Header `X-User-IsAdmin` harus bernilai `1` atau `2`"}
	}

	return c.Next()
}
