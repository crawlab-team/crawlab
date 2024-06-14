package entity

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/trace"
)

type GrpcBaseServiceMessage struct {
	ModelId interfaces.ModelId `json:"id"`
	Data    []byte             `json:"d"`
}

func (msg *GrpcBaseServiceMessage) GetModelId() interfaces.ModelId {
	return msg.ModelId
}

func (msg *GrpcBaseServiceMessage) GetData() []byte {
	return msg.Data
}

func (msg *GrpcBaseServiceMessage) ToBytes() (data []byte) {
	data, err := json.Marshal(*msg)
	if err != nil {
		_ = trace.TraceError(err)
		return data
	}
	return data
}
