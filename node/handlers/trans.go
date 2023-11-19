package handlers

import (
	"chopcoin/shared"
	"chopcoin/shared/models"

	"github.com/gofiber/fiber/v2"
)

func MakeTransaction(ctx *fiber.Ctx) error {
	var input models.SendReq
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}
	balance, err := BC.Balance(input.From)
	if err != nil {
		return err
	}
	if balance < input.Amount {
		return fiber.ErrNotAcceptable
	}
	tx := models.Transaction{
		Sender:    []shared.PublicKey{input.From},
		Recipient: input.To,
		Amount:    input.Amount,
	}
	return ctx.JSON(tx)
}
