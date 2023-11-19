package models

import (
	"chopcoin/shared"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

const genesisCoinbaseData = "The novy god is fucking coming"

// Blockchain implements interactions with a DB
type Blockchain struct {
	tip Block
	Dir string
}

func (bc *Blockchain) makeBlockPath(block Block) string {
	return path.Join(bc.Dir, fmt.Sprint(block.Height))
}

// CreateBlockchain creates a new blockchain DB
func (bc *Blockchain) Create(recipient shared.PublicKey) error {
	if _, err := os.Stat(bc.Dir); err == nil {
		return errors.New("exists")
	}
	if err := os.MkdirAll(bc.Dir, 0700); err != nil {
		return err
	}

	cbtx := NewCoinbaseTX(recipient)
	genesis := NewGenesisBlock(cbtx)
	_, err := bc.MineBlock(genesis.Transactions[0], recipient)
	if err != nil {
		return err
	}

	return nil
}
func (bc *Blockchain) Connect(recipient shared.PublicKey) error {
	fs := os.DirFS(bc.Dir)
	dir, err := os.ReadDir(bc.Dir)
	if err != nil {
		return err
	}
	var highest Block
	for _, b := range dir {
		file, err := fs.Open(b.Name())
		if err != nil {
			return err
		}
		var block Block
		if err := json.NewDecoder(file).Decode(&block); err != nil {
			return err
		}
		if block.Height > highest.Height {
			highest = block
		}
	}
	bc.tip = highest
	fmt.Println(bc.tip.Height)
	return nil
}

func (bc *Blockchain) ContainsBlock(block *Block) bool {
	if _, err := os.Stat(bc.makeBlockPath(*block)); err == nil {
		return true
	}
	return false
}

// AddBlock saves the block into the blockchain
func (bc *Blockchain) AddBlock(block *Block) error {

	if bc.ContainsBlock(block) {
		return errors.New("block already in chain")
	}
	bytes, err := json.Marshal(*block)
	if err != nil {
		return err
	}

	if err := os.WriteFile(bc.makeBlockPath(*block), bytes, 0600); err != nil {
		return err
	}
	bc.tip = *block
	return nil
}

// FindTransaction finds a transaction by its ID

// Iterator returns a BlockchainIterator
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc, nil}
	return bci
}

// GetBlock finds a block by its hash and returns it
func (bc *Blockchain) GetBlock(blockHash shared.Bytes) (Block, error) {
	var block Block

	bytes, err := os.ReadFile(bc.makeBlockPath(Block{Hash: blockHash}))
	if err != nil {
		return block, err
	}
	if err := json.Unmarshal(bytes, &block); err != nil {
		return block, err
	}
	return block, nil
}

// GetBlockHashes returns a list of hashes of all the blocks in the chain
func (bc *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()

	for {
		block := bci.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}
func (bc *Blockchain) Balance(wallet shared.PublicKey) (float64, error) {
	iter := bc.Iterator()
	balance := 0.0
	for {
		block := iter.Next()
		if iter.Error != nil {
			return balance, iter.Error
		}
		for _, tr := range block.Transactions {
			if tr.Recipient.Equals(wallet) {
				balance += tr.Amount
			}
			if len(tr.Sender) != 0 && tr.Sender[0].Equals(wallet) {
				balance -= tr.Amount
			}
		}
		if len(block.PrevBlockHash.String()) == 0 || block.Height == 1 {
			break
		}
	}
	return balance, nil
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transaction *SignedTransaction, recipient shared.PublicKey) (*Block, error) {

	newBlock := NewBlock(transaction, bc.tip.Hash, bc.tip.Height+1)
	newBlock.Transactions = append(newBlock.Transactions,
		&SignedTransaction{Transaction: *NewCoinbaseTX(recipient)})

	data, err := json.Marshal(newBlock)
	if err != nil {
		return newBlock, err
	}
	sum := sha512.Sum512(data)
	newBlock.Hash = sum[:]
	if err := bc.AddBlock(newBlock); err != nil {
		return nil, err
	}

	return newBlock, nil
}
