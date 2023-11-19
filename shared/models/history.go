package models

import (
	"chopcoin/shared"
)

type HistoryEntry struct {
	Who    []shared.PublicKey `json:"who"`
	Amount float64            `json:"amount"`
	Type   string             `json:"type"`
}
