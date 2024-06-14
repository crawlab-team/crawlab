package utils

import "encoding/json"

func JsonToBytes(d interface{}) (bytes []byte, err error) {
	switch d.(type) {
	case []byte:
		return d.([]byte), nil
	default:
		return json.Marshal(d)
	}
}
