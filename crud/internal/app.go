package internal

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/solsteace/goody/crud/internal/controller"
	"github.com/solsteace/goody/crud/internal/repository"
	"github.com/solsteace/goody/crud/internal/route"
	"github.com/solsteace/goody/crud/internal/service"
	"github.com/solsteace/goody/crud/internal/util/api"
	"github.com/solsteace/goody/crud/internal/util/crypto"
	"github.com/solsteace/goody/crud/internal/util/middleware"
	"github.com/solsteace/goody/crud/internal/util/storage"
	"github.com/solsteace/goody/crud/internal/util/token"
	"github.com/solsteace/goody/crud/internal/util/view"
	"github.com/solsteace/goody/lib/payload"
	goodyToken "github.com/solsteace/goody/lib/token"
)

func RunApp() {
	loadEnv()
	upSince := time.Now().Unix()
	db := storage.NewGorm(EnvDbUrl)
	cryptor := crypto.NewBcrypt(10)
	jwtAuth := token.NewJwt[goodyToken.Auth](
		EnvTokenIssuer,
		EnvTokenSecret,
		time.Duration(EnvTokenLifetime))
	indoApi := api.NewEmsifa(EnvIndoApiEndpoint)
	payloader := payload.Rakamin{}
	viewer := view.NewRakamin(indoApi)
	authToken := middleware.NewAuthToken(jwtAuth)
	errorHandler := middleware.NewErrorHandler(payloader)

	// ========================================
	// Layers...
	// ========================================
	alamatRepo := repository.NewGormAlamat(db)
	userRepo := repository.NewGormUser(db)
	authService := service.NewAuth(userRepo, cryptor, jwtAuth)
	alamatService := service.NewAlamat(alamatRepo)
	userService := service.NewUser(userRepo, cryptor)
	authController := controller.NewAuth(&authService, viewer, payloader)
	alamatController := controller.NewAlamat(&alamatService, viewer, payloader)
	userController := controller.NewUser(&userService, viewer, payloader)

	// ========================================
	// Routings...
	// ========================================
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler.Handle,
	})
	app.Use(logger.New())
	api := app.Group("/api")
	route.UseAuth(&api, &authController, authToken)
	route.UseUser(&api, &userController, &alamatController, authToken)
	api.Get("/health", func(c *fiber.Ctx) error {
		upTime := time.Now().Unix() - upSince
		return c.SendString(fmt.Sprintf("%d", upTime))
	})

	app.Listen(fmt.Sprintf(":%d", EnvPort))
}
