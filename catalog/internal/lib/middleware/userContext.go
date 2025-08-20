package middleware

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func UserContext(c *fiber.Ctx) error {
	userId, err := strconv.ParseUint(c.Get("X-User-Id"), 10, 64)
	if err != nil {
		return errors.New("`X-User-Id` should be a positive number")
	}
	c.Locals("userId", uint(userId))

	if isAdmin := c.Get("X-User-IsAdmin"); isAdmin == "0" {
		c.Locals("isAdmin", false)
	} else if isAdmin == "1" {
		c.Locals("isAdmin", true)
	} else {
		return errors.New("`X-User-IsAdmin` value should be either '1' or '2'")
	}

	return c.Next()
}
