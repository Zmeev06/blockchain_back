package handlers

import (
	"chopcoin/shared/models"

	"github.com/gofiber/fiber/v2"
)

func Sync(ctx *fiber.Ctx) error {
	var block models.Block
	if err := ctx.BodyParser(&block); err != nil {
		return err
	}
	return BC.AddBlock(&block)
}
