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
	bc.tip = *genesis
	return nil
}
func (bc *Blockchain) Connect(recipient shared.PublicKey) error {
	dir, err := os.ReadDir(bc.Dir)
	if err != nil {
		return err
	}
	var highest Block
	for _, b := range dir {
		block, err := bc.GetBlock(shared.Bytes(path.Base(b.Name())))
		if err != nil {
			return err
		}
		var block Block
		if err :=json.NewDecoder(file).Decode(&block); err != nil {
		// if err := json.Unmarshal(bytes, &block); err != nil {
			return err
		}
		if block.Height > highest.Height {
			highest = block
		}
	}
	bc.tip = highest
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
		return nil
	}

	if err := os.WriteFile(bc.makeBlockPath(*block), bytes, 0600); err != nil {
		return err
	}
	return nil
}

// FindTransaction finds a transaction by its ID
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}
func (bc *Blockchain) FindSpendableOutputs(recepient shared.PublicKey, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(recepient)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(recepient) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}
func (bc *Blockchain) FindUnspentTransactions(address shared.PublicKey) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.IsLockedWithKey(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Blockchain) FindUTXO() map[string]models.TXOutputs {
	UTXO := make(map[string]models.TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// Iterator returns a BlockchainIterat
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{&bc.tip, bc, nil}
	return bci
}

// GetBestHeight returns the height of the latest block
func (bc *Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

// GetBlock finds a block by its hash and returns it
func (bc *Blockchain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		blockData := b.Get(blockHash)

		if blockData == nil {
			return errors.New("Block is not found.")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
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
