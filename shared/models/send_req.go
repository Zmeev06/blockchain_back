package models

import "chopcoin/shared"

type SendReq struct {
	From   shared.PublicKey `json:"from"`
	To     shared.PublicKey `json:"to"`
	Amount float64          `json:"amount"`
}
