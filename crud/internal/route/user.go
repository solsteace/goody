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

	v1 := user.Group("/v1")
	v1.Get("/", userController.GetProfile)
	v1.Put("/", userController.UpdateProfile)
	v1.Patch("/credentials", userController.ChangeCredentials)
}
