package utils

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/interfaces"
)

func GetResultHash(value interface{}, keys []string) (res string, err error) {
	m := make(map[string]interface{})
	for _, k := range keys {
		_value, ok := value.(interfaces.Result)
		if !ok {
			continue
		}
		v := _value.GetValue(k)
		m[k] = v
	}
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return EncryptMd5(string(data)), nil
}
