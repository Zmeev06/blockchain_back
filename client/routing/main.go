package routing

import (
	"chopcoin/client/handlers"
	"chopcoin/client/middleware"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func Setup(app *fiber.App) {
	// handlers.Init(ctx)
	handlers.Init()
	api := app.Group("/api")
	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)
	api.Use(middleware.Protected(handlers.JWT_SECRET))
	api.Post("/send", handlers.Send)
	api.Get("/user_data", handlers.UserData)
	api.Get("/balance", handlers.Balance)
	api.Get("/history", handlers.History)
	api.Get("/crap", handlers.CrapGet)
	api.Post("/crap", handlers.CrapPost)
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.Dir("dist"),
		NotFoundFile: "index.html"}))
}
