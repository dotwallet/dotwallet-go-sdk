package dotwallet

import "encoding/json"

// ToJSONStr will convert json to string
func ToJSONStr(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(b)
}
