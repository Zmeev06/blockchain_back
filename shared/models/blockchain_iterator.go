package models

import (
	"encoding/json"
	"os"
)

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	nextBlock Block
	bc        *Blockchain
	Error     error
}

// Next returns next block starting from the tip
func (this *BlockchainIterator) Next() Block {

	var block = this.nextBlock
	bytes, err := os.ReadFile(this.bc.makeBlockPath(block))
	if err != nil {
		this.Error = err
		return this.nextBlock
	}
	if err := json.Unmarshal(bytes, &block); err != nil {
		this.Error = err
	}

	this.nextBlock.Height = block.Height - 1

	return block
}
