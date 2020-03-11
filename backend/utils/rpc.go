package utils

import "encoding/json"

// Object 转化为 String
func ObjectToString(params interface{}) string {
	bytes, _ := json.Marshal(params)
	return BytesToString(bytes)
}

// 获取 RPC 参数
func GetRpcParam(key string, params map[string]string) string {
	return params[key]
}
