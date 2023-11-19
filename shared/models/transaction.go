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
}

// Transaction represents a Bitcoin transaction
type Transaction struct {
	Sender    []shared.PublicKey
	Recipient shared.PublicKey
	Amount    float64
}

// Sign signs each input of a Transaction
func (tx *Transaction) Sign(privKey rsa.PrivateKey) (st SignedTransaction, err error) {
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
	// st. = shared.PublicKey(privKey.PublicKey)
	// fmt.Printf("%#v", st)
	// st.ID = hashed[:]
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
