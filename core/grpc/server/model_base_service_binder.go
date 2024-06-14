package server

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
)

func NewModelBaseServiceBinder(req *grpc.Request) (b *ModelBaseServiceBinder) {
	return &ModelBaseServiceBinder{
		req: req,
		msg: &entity.GrpcBaseServiceMessage{},
	}
}

type ModelBaseServiceBinder struct {
	req *grpc.Request
	msg interfaces.GrpcModelBaseServiceMessage
}

func (b *ModelBaseServiceBinder) Bind() (res *entity.GrpcBaseServiceParams, err error) {
	if err := b.bindBaseServiceMessage(); err != nil {
		return nil, err
	}
	params := &entity.GrpcBaseServiceParams{}
	return b.process(params)
}

func (b *ModelBaseServiceBinder) MustBind() (res interface{}) {
	res, err := b.Bind()
	if err != nil {
		panic(err)
	}
	return res
}

func (b *ModelBaseServiceBinder) BindWithBaseServiceMessage() (params *entity.GrpcBaseServiceParams, msg interfaces.GrpcModelBaseServiceMessage, err error) {
	if err := json.Unmarshal(b.req.Data, b.msg); err != nil {
		return nil, nil, err
	}
	params, err = b.Bind()
	if err != nil {
		return nil, nil, err
	}
	return params, b.msg, nil
}

func (b *ModelBaseServiceBinder) process(params *entity.GrpcBaseServiceParams) (res *entity.GrpcBaseServiceParams, err error) {
	if err := json.Unmarshal(b.msg.GetData(), params); err != nil {
		return nil, trace.TraceError(err)
	}
	return params, nil
}

func (b *ModelBaseServiceBinder) bindBaseServiceMessage() (err error) {
	return json.Unmarshal(b.req.Data, b.msg)
}
