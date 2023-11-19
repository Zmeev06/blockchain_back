package routing

import (
	"chopcoin/node/handlers"
	"chopcoin/shared/models"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, bc models.Blockchain) {
	handlers.Init(bc)
	api := app.Group("/api/node")
	api.Get("/mine", handlers.Mine)
	api.Post("/make_trans", handlers.MakeTransaction)
	api.Post("/push_trans", handlers.PushTransaction)
	api.Post("/sync", handlers.Sync)
	api.Post("/balance", handlers.Balance)
	api.Post("/history", handlers.History)
}
