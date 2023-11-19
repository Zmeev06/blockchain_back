package shared

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
)

type PublicKey rsa.PublicKey


func (this PublicKey) MarshalJSON() ([]byte, error) {
	v := hex.EncodeToString(x509.MarshalPKCS1PublicKey((*rsa.PublicKey)(&this)))
	return json.Marshal(v)
}
func (this *PublicKey) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	bytes, err := hex.DecodeString(v)
	if err != nil {
		return err
	}
	key, err := x509.ParsePKCS1PublicKey(bytes)
	if err != nil {
		return err
	}
	this = (*PublicKey)(key)
	return err
}
