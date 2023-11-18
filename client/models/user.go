package models

import (
	"crypto/rsa"
)

type User struct {
	Login    string           `json:"login"`
	Password string           `json:"password"`
	Privkey  rsa.PrivateKey   `json:"priv_key"`
}
