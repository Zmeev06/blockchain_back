package handlers

import (
	"chopcoin/shared"
	"chopcoin/shared/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Balance(ctx *fiber.Ctx) error {
	user, err := getUserFromJwt(ctx)
	if err != nil {
		return err
	}
	data := models.BalanceReq{
		PubKey: shared.PublicKey(user.Privkey.PublicKey),
	}
	agent := fiber.Post(fmt.Sprintf("http://localhost%s/api/node/balance", NODE_ADDR))
	status, bytes, errs := agent.JSON(data).Bytes()
	ctx.Context().SetBody(bytes)
	if len(errs) != 0 {
		return ctx.Status(status).JSON(errs)
	}
	if status != 200 {
		return ctx.SendStatus(status)
	}
	return nil
}
