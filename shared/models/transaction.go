package models

import (
	"chopcoin/shared"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"

	"encoding/json"
)

const subsidy = 10

type SignedTransaction struct {
	Transaction
	Signature shared.Bytes
	PubKey    shared.PublicKey
	ID        shared.Bytes
}

// Transaction represents a Bitcoin transaction
type Transaction struct {
	Sender []shared.PublicKey
	Recipient shared.PublicKey
	Amount float64
}

// Sign signs each input of a Transaction
func (tx *Transaction) Sign(privKey rsa.PrivateKey, prevTXs map[string]Transaction) (st SignedTransaction, err error) {
	data, err := json.Marshal(*tx)
	if err != nil {
		return
	}
	hashed := sha512.Sum512(data)
	sig, err := rsa.SignPKCS1v15(rand.Reader, &privKey, crypto.SHA512, hashed[:])
	if err != nil {
		return
	}
	st.Transaction = *tx
	st.Signature = sig
	st.PubKey = shared.PublicKey(privKey.PublicKey)
	st.ID = hashed[:]
	return
}

// Verify verifies signatures of Transaction inputs
func (tx *SignedTransaction) Verify() bool {

	if len(tx.Sender) == 0 {
		return true
	}
	data, err := json.Marshal(tx.Transaction)
	if err != nil {
		return false
	}

	if rsa.VerifyPKCS1v15((*rsa.PublicKey)(&tx.Sender[0]), crypto.SHA512, data, tx.Signature) != nil {
		return false
	}

	return true
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to shared.PublicKey) *Transaction {

	tx := Transaction{[]shared.PublicKey{}, to, subsidy}

	return &tx
}

// NewUTXOTransaction creates a new transaction
func NewUTXOTransaction(wallet *rsa.PublicKey, to *rsa.PublicKey, amount int, UTXOSet *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	from := fmt.Sprintf("%s", wallet.GetAddress())
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXOSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}

// DeserializeTransaction deserializes a transaction
func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}
