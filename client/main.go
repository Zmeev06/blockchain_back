package main

import (
	"chopcoin/client/routing"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	routing.Setup(app)
	log.Fatal(app.Listen(os.Getenv("LISTEN_ADDR")))
}
