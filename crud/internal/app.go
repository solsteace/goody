package internal

import (
	"fmt"
	"log"
	"os"
	"path"
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

// Preparing save paths
func RunApp() {
	// ========================================
	// Initializations...
	// ========================================
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
	// userContext := middleware.NewUserContext()

	savePaths := map[string]string{
		"toko":   path.Join(EnvUploadSaveBasePath, "toko"),
		"produk": path.Join(EnvUploadSaveBasePath, "produk")}
	for _, savePath := range savePaths {
		if err := os.MkdirAll(savePath, 0777); err != nil {
			log.Fatalf("Failed creating directory: %v", err)
		}
	}

	// ========================================
	// Layers...
	// ========================================
	alamatRepo := repository.NewGormAlamat(db)
	userRepo := repository.NewGormUser(db)
	// kategoriRepo := repository.NewGormKategori(db)
	// produkRepo := repository.NewGormProduk(db)
	// tokoRepo := repository.NewGormToko(db)

	authService := service.NewAuth(userRepo, cryptor, jwtAuth)
	alamatService := service.NewAlamat(alamatRepo)
	userService := service.NewUser(userRepo, cryptor)
	// kategoriService := service.NewKategori(kategoriRepo)
	// tokoService := service.NewToko(tokoRepo, savePaths["toko"])
	// produkService := service.NewProduk(produkRepo)

	authController := controller.NewAuth(&authService, viewer, payloader)
	alamatController := controller.NewAlamat(&alamatService, viewer, payloader)
	userController := controller.NewUser(&userService, viewer, payloader)
	// kategoriController := controller.NewKategori(kategoriService, viewer, payloader)
	// tokoControler := controller.NewToko(tokoService, viewer)
	// produkController := controller.NewProduk(produkService, viewer, payloader)

	// ========================================
	// Routings...
	// ========================================
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler.Handle,
	})
	app.Use(logger.New())
	api := app.Group("/api")
	v1 := api.Group("/v1")
	route.UseAuth(&v1, &authController, authToken)
	route.UseUser(&v1, &userController, &alamatController, authToken)
	// route.UseProduk(&v1, &produkController, authToken, userContext)
	// route.UseToko(&v1, &tokoControler)
	// route.UseKategori(&v1, &kategoriController)
	api.Get("/health", func(c *fiber.Ctx) error {
		upTime := time.Now().Unix() - upSince
		return c.SendString(fmt.Sprintf("%d", upTime))
	})

	app.Listen(fmt.Sprintf(":%d", EnvPort))
}
