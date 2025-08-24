package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/controller"
	"github.com/solsteace/goody/account/internal/lib/middleware"
)

func UseUser(
	parent *fiber.Router,
	userController *controller.User,
	alamatController *controller.Alamat,
	authToken middleware.AuthToken,
) {
	user := (*parent).Group("/user")
	user.Use(authToken.Handle)

	UseAlamat(&user, alamatController)
	userV1 := user.Group("/v1")
	userV1.Get("/", userController.GetProfile)
	userV1.Put("/", userController.UpdateProfile)
	userV1.Patch("/credentials", userController.ChangeCredentials)
}
