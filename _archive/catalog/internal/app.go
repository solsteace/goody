package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"github.com/solsteace/goody/catalog/internal/controller"
	"github.com/solsteace/goody/catalog/internal/lib/persistence"
	"github.com/solsteace/goody/catalog/internal/lib/view"
	"github.com/solsteace/goody/catalog/internal/repository"
	"github.com/solsteace/goody/catalog/internal/route"
	"github.com/solsteace/goody/catalog/internal/service"
	"github.com/solsteace/goody/lib/messaging/event"
	"github.com/solsteace/goody/lib/payload"
)

func RunApp() {
	loadEnv()
	upSince := time.Now().Unix()
	db := persistence.NewGorm(EnvDbUrl)
	savePaths := map[string]string{
		"toko":   path.Join(EnvUploadSaveBasePath, "toko"),
		"produk": path.Join(EnvUploadSaveBasePath, "produk")}
	viewer := view.NewRakamin()
	payloader := payload.Rakamin{}

	// Layers...
	kategoriRepo := repository.NewGormKategori(db)
	produkRepo := repository.NewGormProduk(db)
	tokoRepo := repository.NewGormToko(db)
	kategoriService := service.NewKategori(kategoriRepo)
	tokoService := service.NewToko(tokoRepo, savePaths["toko"])
	produkService := service.NewProduk(produkRepo)
	_ = controller.NewKategori(kategoriService, viewer, payloader)
	tokoController := controller.NewToko(tokoService, viewer)
	produkController := controller.NewProduk(produkService)

	// Routes...
	app := fiber.New()
	api := app.Group("/api")
	route.UseProduk(&api, &produkController)
	route.UseToko(&api, &tokoController)
	api.Get("/health", func(c *fiber.Ctx) error {
		upTime := time.Now().Unix() - upSince
		return c.SendString(fmt.Sprintf("%d", upTime))
	})

	// Subscriptions, side effects...
	// Preparing save paths
	for _, savePath := range savePaths {
		if err := os.MkdirAll(savePath, 0777); err != nil {
			log.Fatalf("Failed creating directory: %v", err)
		}
	}

	// MQ setup
	// TODO: Refactor and add auto-reconnection
	mqConn, err := amqp091.Dial(EnvMqUrl)
	if err != nil {
		log.Fatalf("Couldn't connect to MQ: %v", err)
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

	msg, err := channel.Consume(
		queue.Name,
		"create.toko",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatalf("Couldn't consume messages: %v", err)
	}
	go (func() {
		for m := range msg {
			payload := new(event.UserRegistered)
			if err := json.Unmarshal(m.Body, payload); err != nil {
				fmt.Printf("Warning! invalid message: %v \n", err)
				m.Nack(false, false)
			}

			if _, err := tokoService.Create(payload.IdUser, payload.Nama); err != nil {
				fmt.Printf("Warning! failure during creating toko: %v\n", err)
				m.Nack(false, true)
			}

			fmt.Printf("Successfully create `%s` toko", payload.Nama)
			m.Ack(false)
		}
	})()

	app.Listen(fmt.Sprintf(":%d", EnvPort))
}
