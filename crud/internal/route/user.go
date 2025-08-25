package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/crud/internal/controller"
	"github.com/solsteace/goody/crud/internal/util/middleware"
)

func UseUser(
	parent *fiber.Router,
	userController *controller.User,
	alamatController *controller.Alamat,
	authToken middleware.AuthChecker,
) {
	user := (*parent).Group("/user")

	user.Use(authToken.Handle)
	UseAlamat(&user, alamatController)

	user.Get("/", userController.GetProfile)
	user.Put("/", userController.UpdateProfile)
	user.Patch("/credentials", userController.ChangeCredentials)
}
