package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rabbitmq/amqp091-go"
	"github.com/solsteace/goody/account/internal/controller"
	"github.com/solsteace/goody/account/internal/domain"
	"github.com/solsteace/goody/account/internal/lib/api"
	"github.com/solsteace/goody/account/internal/lib/crypto"
	"github.com/solsteace/goody/account/internal/lib/messaging/amqp"
	"github.com/solsteace/goody/account/internal/lib/middleware"
	"github.com/solsteace/goody/account/internal/lib/persistence"
	"github.com/solsteace/goody/account/internal/lib/token"
	"github.com/solsteace/goody/account/internal/repository"
	"github.com/solsteace/goody/account/internal/route"
	"github.com/solsteace/goody/account/internal/service"
	"github.com/solsteace/goody/lib/token/payload"
)

func RunApp() {
	loadEnv()
	upSince := time.Now().Unix()
	db := persistence.NewGorm(EnvDbUrl)
	cryptor := crypto.NewBcrypt(10)
	jwtAuth := token.NewJwt[payload.AuthPayload](
		EnvTokenIssuer,
		EnvTokenSecret,
		time.Duration(EnvTokenLifetime))
	indoApi := api.NewEmsifa(EnvIndoApiEndpoint)
	authToken := middleware.NewAuthToken(jwtAuth)

	// Layers...
	alamatRepo := repository.NewGormAlamat(db)
	userRepo := repository.NewGormUser(db)
	authService := service.NewAuth(userRepo, cryptor, indoApi, jwtAuth)
	alamatService := service.NewAlamat(alamatRepo)
	userService := service.NewUser(userRepo, cryptor, indoApi)
	authController := controller.NewAuth(&authService)
	alamatController := controller.NewAlamat(&alamatService)
	userController := controller.NewUser(&userService)

	// Routings...
	app := fiber.New()
	app.Use(logger.New())
	api := app.Group("/api")
	route.UseAuth(&api, &authController, authToken)
	route.UseUser(&api, &userController, &alamatController, authToken)
	api.Get("/health", func(c *fiber.Ctx) error {
		upTime := time.Now().Unix() - upSince
		return c.SendString(fmt.Sprintf("%d", upTime))
	})

	// Subscriptions, side-effects...
	queue := amqp.NewQueue("hello", false, false, false, false, nil)
	channel := amqp.NewChannel()
	channel.Track(&queue)
	publisher := amqp.NewPublisher(EnvMqUrl, 2)
	publisher.Track(&channel, "myChannel")
	publisher.Initiate()

	authService.SubscribeOnNewUser(func(u domain.User) error {
		body, err := json.Marshal(struct {
			IdUser uint   `json:"id_user"`
			Nama   string `json:"nama"`
		}{IdUser: u.ID, Nama: u.Nama})
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Q: No binding between an exchange and the queue?
		// A: https://www.cloudamqp.com/blog/part4-rabbitmq-for-beginners-exchanges-routing-keys-bindings.html#default-exchange
		channel.Conn.PublishWithContext(
			ctx,
			"",              // exchange
			queue.Conn.Name, // routing key
			false,           // mandatory
			false,           // immediate
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			return err
		}
		return nil
	})

	app.Listen(fmt.Sprintf(":%d", EnvPort))
}
