package main

import (
	"github.com/solsteace/goody/account/internal"
)

func main() {
	app := internal.NewApp()
	app.Listen(":8880")
}
