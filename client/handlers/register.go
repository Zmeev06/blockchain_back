package handlers

import (
	"chopcoin/client/models"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"

	"os"

	"github.com/gofiber/fiber/v2"
)

func Register(ctx *fiber.Ctx) error {
	var input creds
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}
	userPath := makeUserPath(input.Login)
	if _, err := os.Stat(userPath); err == nil {
		return fiber.ErrForbidden
	}
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	userModel := models.User{
		Login:    input.Login,
		Password: input.Password,
		Privkey:  *priv,
	}
	user, err := json.Marshal(userModel)
	if err != nil {
		return err
	}
	if err := os.WriteFile(userPath, user, 0600); err != nil {
		return err
	}
	str, err := makeToken(userModel)
	ctx.WriteString(str)
	return nil
}
