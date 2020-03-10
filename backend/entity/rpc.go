package entity

type RpcMessage struct {
	Id      string            `json:"id"`
	Method  string            `json:"method"`
	NodeId  string            `json:"node_id"`
	Params  map[string]string `json:"params"`
	Timeout int               `json:"timeout"`
	Result  string            `json:"result"`
	Error   string            `json:"error"`
}
