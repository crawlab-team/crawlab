package utils

import "encoding/json"

func GetObjectHash(obj any) string {
	data, _ := json.Marshal(obj)
	if data == nil {
		// random hash
		return EncryptMd5(NewUUIDString())
	}
	return EncryptMd5(string(data))
}
