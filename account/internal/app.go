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
	"github.com/solsteace/goody/account/internal/repository"
	"github.com/solsteace/goody/account/internal/route"
	"github.com/solsteace/goody/account/internal/service"
	"github.com/solsteace/goody/lib/token/payload"
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

	alamatRepo := repository.NewGormAlamat(db)
	userRepo := repository.NewGormUser(db)
	authService := service.NewAuth(userRepo, cryptor, indoApi, jwtAuth)
	alamatService := service.NewAlamat(alamatRepo)
	userService := service.NewUser(userRepo, cryptor, indoApi)
	authController := controller.NewAuth(authService)
	alamatController := controller.NewAlamat(alamatService)
	userController := controller.NewUser(userService)

	app := fiber.New()
	api := app.Group("/api")
	route.UseAuth(&api, &authController)
	route.UseUser(&api, &userController, &alamatController, &jwtAuth)
	api.Get("/health", func(c *fiber.Ctx) error {
		upTime := time.Now().Unix() - upSince
		return c.SendString(fmt.Sprintf("%d", upTime))
	})

	return app
}
