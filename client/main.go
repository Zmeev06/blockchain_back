package main

import (
	"chopcoin/client/routing"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	routing.Setup(app)
	log.Fatal(app.Listen(os.Getenv("LISTEN_ADDR")))
}
