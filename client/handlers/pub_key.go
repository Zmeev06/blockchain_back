package handlers

import (
	"chopcoin/shared"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func UserData(ctx *fiber.Ctx) error {
	user, err := getUserFromJwt(ctx)
	if err != nil {
		return err
	}
	return ctx.JSON(
		jwt.MapClaims{
			"pub_key": shared.PublicKey(user.Privkey.PublicKey),
			"login":   user.Login,
		},
	)
}
