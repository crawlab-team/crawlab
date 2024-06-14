package entity

type RpcMessage struct {
	Id      string            `json:"id"`      // 消息ID
	Method  string            `json:"method"`  // 消息方法
	NodeId  string            `json:"node_id"` // 节点ID
	Params  map[string]string `json:"params"`  // 参数
	Timeout int               `json:"timeout"` // 超时
	Result  string            `json:"result"`  // 结果
	Error   string            `json:"error"`   // 错误
}
