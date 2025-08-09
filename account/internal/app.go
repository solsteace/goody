package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/solsteace/goody/account/internal/lib/persistence"
	"github.com/solsteace/goody/account/internal/repository"
)

func NewApp() *fiber.App {
	upSince := time.Now().Unix()
	db := persistence.NewGorm(os.Getenv("DB_URL"))

	// Prepare layers...
	userRepo := repository.NewGormUserRepo(db)
	authService := NewAuthService(userRepo)
	authController := NewAuthController(authService)

	// Prepare endpoints...
	app := fiber.New()
	api := app.Group("/api")
	registerAuthRoutes(&api, &authController)

	api.Get("/health", func(c *fiber.Ctx) error {
		upTime := time.Now().Unix() - upSince
		return c.SendString(fmt.Sprintf("%d", upTime))
	})

	// Prepare one-off calls and routines
	userRepo.Migrate()

	return app
}
