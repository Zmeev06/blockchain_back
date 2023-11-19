package models

import (
	"bytes"
	"chopcoin/shared"
	"encoding/gob"
	"log"
)

// TXOutput represents a transaction output
type TXOutput struct {
	Value  int
	PubKey shared.PublicKey
}

// Lock signs the output
func (out *TXOutput) Lock(address shared.PublicKey) {
	out.PubKey = address
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (out *TXOutput) IsLockedWithKey(pubKeyHash shared.PublicKey) bool {
	return bytes.Compare(out.PubKey, pubKeyHash) == 0
}

// NewTXOutput create a new TXOutput
func NewTXOutput(value int, recepient shared.PublicKey) *TXOutput {
	txo := &TXOutput{value, recepient}
	return txo
}

// TXOutputs collects TXOutput
type TXOutputs struct {
	Outputs []TXOutput
}

// Serialize serializes TXOutputs
func (outs TXOutputs) Serialize() ([]byte, error) {
	return  json.Marshal(outs)
}

// DeserializeOutputs deserializes TXOutputs
func DeserializeOutputs(data []byte) (TXOutputs, error) {
	var outputs TXOutputs

	err := json.Unmarshal(data, &outputs)
	if err != nil {
		return outputs, err
	}

	return outputs, nil
}
