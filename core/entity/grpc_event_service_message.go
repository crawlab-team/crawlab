package entity

type GrpcEventServiceMessage struct {
	Type   string   `json:"type"`
	Events []string `json:"events"`
	Key    string   `json:"key"`
	Data   []byte   `json:"data"`
}
