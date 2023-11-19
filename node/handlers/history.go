package handlers

import (
	"chopcoin/shared"
	"chopcoin/shared/models"

	"github.com/gofiber/fiber/v2"
)

func History(ctx *fiber.Ctx) error {
	var input models.BalanceReq
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}
	wallet := input.PubKey
	iter := BC.Iterator()
	items := []models.HistoryEntry{}
	for {
		block := iter.Next()
		if iter.Error != nil {
			return iter.Error
		}
		for _, tr := range block.Transactions {
			if tr.Recipient.Equals(wallet) {
				item := models.HistoryEntry{
					Who:    tr.Sender,
					Amount: tr.Amount,
				}
				if len(tr.Sender) == 0 {
					item.Type = "refill"
				} else {
					item.Type = "incoming"
				}
				items = append(items, item)
			}
			if len(tr.Sender) != 0 && tr.Sender[0].Equals(wallet) {
				items = append(items, models.HistoryEntry{
					Who:    []shared.PublicKey{tr.Recipient},
					Amount: tr.Amount,
					Type:   "outgoing",
				})
			}
		}
		if block.Height == 1 || len(block.PrevBlockHash.String()) == 0 {
			break
		}
	}
	return ctx.JSON(items)
}
