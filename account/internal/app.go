package internal

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/controller"
	"github.com/solsteace/goody/account/internal/lib/api"
	"github.com/solsteace/goody/account/internal/lib/crypto"
	"github.com/solsteace/goody/account/internal/lib/persistence"
	"github.com/solsteace/goody/account/internal/lib/token"
	"github.com/solsteace/goody/account/internal/lib/token/payload"
	"github.com/solsteace/goody/account/internal/repository"
	"github.com/solsteace/goody/account/internal/route"
	"github.com/solsteace/goody/account/internal/service"
)

func NewApp() *fiber.App {
	loadEnv()

	upSince := time.Now().Unix()
	db := persistence.NewGorm(EnvDbUrl)
	cryptor := crypto.NewBcrypt(10)
	jwtAuth := token.NewJwt[payload.AuthPayload](
		EnvTokenIssuer,
		EnvTokenSecret,
		time.Duration(EnvTokenLifetime))
	indoApi := api.NewEmsifa(EnvIndoApiEndpoint)

	userRepo := repository.NewGormUserRepo(db)
	authService := service.NewAuthService(userRepo, cryptor, indoApi, jwtAuth)
	userService := service.NewUserService(userRepo, cryptor, indoApi)
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)

	app := fiber.New()
	api := app.Group("/api")
	route.RegisterAuthRoutes(&api, &authController)
	route.RegisterUserRoutes(&api, &userController, &jwtAuth)
	api.Get("/health", func(c *fiber.Ctx) error {
		upTime := time.Now().Unix() - upSince
		return c.SendString(fmt.Sprintf("%d", upTime))
	})

	return app
}
