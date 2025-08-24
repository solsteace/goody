package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/lib/oops"
)

type adminChecker struct{}

func NewAdminChecker() adminChecker {
	return adminChecker{}
}

func (ac adminChecker) Handle(c *fiber.Ctx) error {
	isAdmin, ok := c.Locals("isAdmin").(bool)
	if !ok || !isAdmin {
		return oops.Forbidden{
			Err: errors.New("User should be an admin"),
			Msg: "Anda perlu akses sebagai admin untuk melakukan aksi ini"}
	}

	return c.Next()
}
