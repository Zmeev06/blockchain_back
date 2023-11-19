package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Login(ctx *fiber.Ctx) error {
	var input creds
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}
	user, err := getUserByName(input.Login)
	if err != nil {
		return err
	}
	if input.Password != user.Password {
		return fiber.ErrUnauthorized
	}
	str, err := makeToken(user)
	if err != nil {
		return err
	}
	return ctx.JSON(str)
}
