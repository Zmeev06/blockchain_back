package handlers

import (
	"chopcoin/shared/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func PushTransaction(ctx *fiber.Ctx) error {
	fmt.Println("started push")
	var tx models.SignedTransaction
	if err := ctx.BodyParser(&tx); err != nil {
		return err
	}
	messages<-tx
	return nil
}
