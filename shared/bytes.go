package shared

import (
	"encoding/hex"
	"encoding/json"
)

type Bytes []byte

func (this Bytes) MarshalJSON() ([]byte, error) {
	v := hex.EncodeToString(this)
	return json.Marshal(v)
}
func (this *Bytes) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	bytes, err := hex.DecodeString(v)
	*this = bytes
	return err
}
func (this Bytes) String() string {
	return hex.EncodeToString(this)
}
