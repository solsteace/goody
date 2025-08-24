package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/lib/oops"
	"github.com/solsteace/goody/lib/token"
)

type adminChecker struct{}

func NewAdminChecker() adminChecker {
	return adminChecker{}
}

func (ac adminChecker) Handle(c *fiber.Ctx) error {
	auth, ok := c.Locals("Authorization").(token.Auth)
	if !ok || !auth.IsAdmin {
		return oops.Forbidden{
			Err: errors.New("User should be an admin"),
			Msg: "Anda perlu akses sebagai admin untuk melakukan aksi ini"}
	}

	return c.Next()
}
