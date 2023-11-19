package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func CrapPost(ctx *fiber.Ctx) error {
	user, err := getUserFromJwt(ctx)
	if err != nil {
		return err
	}
	if err := os.WriteFile(makeCrapPath(user.Login), ctx.Body(), 0777); err != nil {
		return err
	}
	return nil
}
func CrapGet(ctx *fiber.Ctx) error {
	user, err := getUserFromJwt(ctx)
	if err != nil {
		return err
	}
	bytes, err := os.ReadFile(makeCrapPath(user.Login)); if err != nil {
		return err
	}
	ctx.Context().SetBody(bytes)
	return nil
}
