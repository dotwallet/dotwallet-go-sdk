package dotwallet

import "encoding/json"

func ToJsonStr(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(b)
}
