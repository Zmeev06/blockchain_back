package models

import (
	"bytes"
	"chopcoin/shared"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
)

// TXInput represents a transaction input
type TXInput struct {
	Txid      shared.Bytes
	Vout      int
	Signature shared.Bytes
	PubKey    shared.PublicKey
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKey shared.PublicKey) bool {
	return in.PubKey == pubKey
}
func (in *TXInput) CanUnlockOutputWith(unlockingData shared.PublicKey) bool {
	return in.PubKey == unlockingData
}
func (in *TXInput) Sign(privKey rsa.PrivateKey, recepient shared.PublicKey) {
	in.PubKey = shared.PublicKey(privKey.PublicKey)
	bts := bytes.Join([][]byte{in.Txid, in.Vout}, []byte{})
	hashed := sha512.Sum512(bts)
	in.Signature = rsa.SignPKCS1v15( rand.Reader, &privKey, crypto.SHA512, hashed[:])
}
