package handlers

import (
	"chopcoin/shared"
	"chopcoin/shared/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func lookupKey(key shared.PublicKey) (user cmodels.User, err error) {
	fs := os.DirFS(USERS)
	dir, err := os.ReadDir(USERS)
	if err != nil {
		return
	}
	for _, v:= range dir {
		file, err := fs.Open(v.Name())
		if err != nil {
			continue
		}
		if err := json.NewDecoder(file).Decode(user); err != nil {
			continue
		}
		if key.Equals(shared.PublicKey(user.Privkey.PublicKey)) {
			return user, nil
		}
	}
	return user, fiber.ErrNotFound
}
func History(ctx *fiber.Ctx) error {
	user, err := getUserFromJwt(ctx)
	if err != nil {
		return err
	}
	data := models.BalanceReq{
		PubKey: shared.PublicKey(user.Privkey.PublicKey),
	}
	agent := fiber.Post(fmt.Sprintf("http://localhost%s/api/node/history", NODE_ADDR))
	status, bytes, errs := agent.JSON(data).Bytes()
	if len(errs) != 0 {
		return ctx.Status(status).JSON(errs)
	}
	var items []models.HistoryEntry
	if err := json.Unmarshal(bytes, &items); err != nil {
		return err
	}
	type HistoryEntryView struct {
		Who    []string  `json:"who"`
		Amount float64 `json:"amount"`
		Type   string  `json:"type"`
	}
	output := make ([]HistoryEntryView, len(items))
	for k, v := range items {
		who := []string{}
		for k, key := range v.Who {
			user, err := lookupKey(key)
			if err != nil {
				return err
			}
			who[k] = user.Login
		}
		output[k] = HistoryEntryView{who, v.Amount, v.Type}
	}

	return ctx.JSON(output)
}
