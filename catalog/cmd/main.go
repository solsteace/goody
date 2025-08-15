package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	// app := internal.NewApp()
	app := fiber.New()
	app.Listen(":8880")
}
