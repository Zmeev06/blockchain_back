package handlers

import (
	"chopcoin/shared"
	"chopcoin/shared/models"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Send(ctx *fiber.Ctx) error {
	type req struct {
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}
	var input req
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}
	sender, err := getUserFromJwt(ctx)
	if err != nil {
		return err
	}
	receiver, err := getUserByName(input.To)
	if err != nil {
		return ctx.SendStatus(404)
	}

	body := models.SendReq{
		From:   shared.PublicKey(sender.Privkey.PublicKey),
		To:     shared.PublicKey(receiver.Privkey.PublicKey),
		Amount: input.Amount,
	}
	addr := NODE_ADDR
	transMake := fiber.Post(fmt.Sprintf("http://localhost%s/api/node/make_trans", addr))
	transMake.JSON(body)
	status, bytes, errs := transMake.Bytes()
	if len(errs) != 0 {
		return ctx.Status(status).JSON(errs)
	}
	if status != 200 {
		return ctx.SendStatus(status)
	}
	var tx models.Transaction
	if err := json.Unmarshal(bytes, &tx); err != nil {
		return err
	}
	stx, err := tx.Sign(sender.Privkey)
	if err != nil {
		return err
	}
	transPush := fiber.Post(fmt.Sprintf("http://localhost%s/api/node/push_trans", addr))
	transPush.JSON(stx)
	status, bytes, errs = transPush.Bytes()
	if len(errs) != 0 {
		return ctx.Status(status).JSON(errs)
	}
	if status != 200 {
		ctx.Context().SetBody(bytes)
		return ctx.SendStatus(status)
	}

	ctx.Context().SetBody(bytes)
	return nil
}
