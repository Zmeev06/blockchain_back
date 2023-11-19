package handlers

import (
	"chopcoin/shared/models"

	"github.com/gofiber/fiber/v2"
)

func Mine(ctx *fiber.Ctx) error {
	in := make(chan models.SignedTransaction)
	entering <- in
	stx := <-in
	leaving <- in
	return ctx.JSON(stx)
}
