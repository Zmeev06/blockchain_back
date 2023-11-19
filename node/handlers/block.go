package handlers

import (
	"chopcoin/node/models"

	"github.com/gofiber/fiber/v2"
)

func AddTransaction(ctx *fiber.Ctx) error {
	var tx models.Transaction
	if err := ctx.BodyParser(&tx); err != nil {
		return err
	}
	for _, ch := range candidates {
		ch<-tx
	}
	return nil
}
