package shared

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
)

func MarshalRsa(key *rsa.PublicKey) string {
	return hex.EncodeToString(x509.MarshalPKCS1PublicKey(key))
}
func UnmarshalRsa(key string) (*rsa.PublicKey, error) {
	bytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PublicKey(bytes)
}
