package handlers

import (
	"chopcoin/shared/models"
)

var BC models.Blockchain

type miner chan<- models.SignedTransaction

var (
	entering = make(chan miner)
	leaving  = make(chan miner)
	messages = make(chan models.SignedTransaction)
)

func Init(bc models.Blockchain) {
	BC = bc
	go broadcaster()
}
func broadcaster() {
	candidates := make(map[miner]bool)
	for {
		select {
		case msg := <-messages:
			for m := range candidates {
				m <- msg
			}
		case cli := <-entering:
			candidates[cli] = true
		case cli := <-leaving:
			delete(candidates, cli)
			close(cli)
		}
	}
}
