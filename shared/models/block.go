package models

import (
	"chopcoin/shared"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	Transactions   []*SignedTransaction
	PrevBlockHash shared.Bytes
	Hash          shared.Bytes
	Nonce         int
	Height        int
}

// NewBlock creates and returns Block
func NewBlock(transaction *SignedTransaction, prevBlockHash []byte, height int) *Block {
	block := &Block{
		time.Now().Unix(),
		[]*SignedTransaction{transaction},
		prevBlockHash,
		[]byte{},
		0,
		height,
	}
	// pow := NewProofOfWork(block)
	// nonce, hash := pow.Run()
	//
	// block.Hash = hash[:]
	// block.Nonce = nonce

	// fmt.Printf("%#v", block)
	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock(
		&SignedTransaction{
			Transaction: *coinbase,
		}, []byte{}, 0)
}
