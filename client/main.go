package main

import (
	"chopcoin/client/routing"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routing.Setup(app)
	log.Fatal(app.Listen(os.Getenv("LISTEN_ADDR")))
}
