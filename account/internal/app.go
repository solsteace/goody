package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rabbitmq/amqp091-go"
	"github.com/solsteace/goody/account/internal/controller"
	"github.com/solsteace/goody/account/internal/domain"
	"github.com/solsteace/goody/account/internal/lib/api"
	"github.com/solsteace/goody/account/internal/lib/crypto"
	"github.com/solsteace/goody/account/internal/lib/middleware"
	"github.com/solsteace/goody/account/internal/lib/persistence"
	"github.com/solsteace/goody/account/internal/lib/token"
	"github.com/solsteace/goody/account/internal/lib/view/rakamin"
	"github.com/solsteace/goody/account/internal/repository"
	"github.com/solsteace/goody/account/internal/route"
	"github.com/solsteace/goody/account/internal/service"
	"github.com/solsteace/goody/lib/messaging/event"
	goodyToken "github.com/solsteace/goody/lib/token"
)

func RunApp() {
	loadEnv()
	upSince := time.Now().Unix()
	db := persistence.NewGorm(EnvDbUrl)
	cryptor := crypto.NewBcrypt(10)
	jwtAuth := token.NewJwt[goodyToken.Auth](
		EnvTokenIssuer,
		EnvTokenSecret,
		time.Duration(EnvTokenLifetime))
	indoApi := api.NewEmsifa(EnvIndoApiEndpoint)
	authToken := middleware.NewAuthToken(jwtAuth)
	viewer := rakamin.NewViewer(indoApi)

	// Layers...
	alamatRepo := repository.NewGormAlamat(db)
	userRepo := repository.NewGormUser(db)
	authService := service.NewAuth(userRepo, cryptor, jwtAuth)
	alamatService := service.NewAlamat(alamatRepo)
	userService := service.NewUser(userRepo, cryptor)
	authController := controller.NewAuth(&authService, viewer)
	alamatController := controller.NewAlamat(&alamatService, viewer)
	userController := controller.NewUser(&userService, viewer)

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
	// TODO: Refactor and add auto-reconnection
	mqConn, err := amqp091.Dial(EnvMqUrl)
	if err != nil {
		log.Fatalf("Couldn't connection to MQ: %v", err)
	}
	defer mqConn.Close()

	channel, err := mqConn.Channel()
	if err != nil {
		log.Fatalf("Couldn't create channel: %v", err)
	}
	defer channel.Close()

	exchangeName := "goody"
	err = channel.ExchangeDeclare(
		exchangeName,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil)   // arguments
	if err != nil {
		log.Fatalf("Couldn't create exchange: %v", err)
	}

	queue, err := channel.QueueDeclare(
		"new.toko", // name
		true,       // durable (on shutdown, should the queue be persisted?)
		false,      // delete when unused (when the last consumer leaves, should the queue be deleted?)
		false,      // exclusive (When the disconnected, should the queue be deleted?)
		false,      // no-wait
		nil)        // arguments
	if err != nil {
		log.Fatalf("Couldn't create queue: %v", err)
	}

	err = channel.QueueBind(
		queue.Name,
		event.UserRegisteredName,
		exchangeName,
		false,
		nil)
	if err != nil {
		log.Fatalf("Couldn't bind `%s` queue to `%s` exchange: %v",
			queue.Name, exchangeName, err)
	}

	authService.SubscribeOnNewUser(func(u domain.User) error {
		payload := event.NewUserRegistered(u.ID, u.Nama)
		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// Q: No binding between an exchange and the queue?
		// A: https://www.cloudamqp.com/blog/part4-rabbitmq-for-beginners-exchanges-routing-keys-bindings.html#default-exchange
		channel.PublishWithContext(
			ctx,
			exchangeName,             // exchange
			event.UserRegisteredName, // routing key
			false,                    // mandatory
			false,                    // immediate
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
