package handlers

import (
	"chopcoin/shared/models"

	"github.com/gofiber/fiber/v2"
)

func Balance(ctx *fiber.Ctx) error {
	var input models.BalanceReq
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}
	balance, err := BC.Balance(input.PubKey)
	if err != nil {
		return err
	}
	return ctx.JSON(balance)
}
