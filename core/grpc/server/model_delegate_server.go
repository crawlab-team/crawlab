package server

import (
	"context"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	grpc "github.com/crawlab-team/crawlab/grpc"
)

type ModelDelegateServer struct {
	grpc.UnimplementedModelDelegateServer
}

// Do and perform an RPC action of constants.Delegate
func (svr ModelDelegateServer) Do(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// bind message
	obj, msg, err := NewModelDelegateBinder(req).BindWithDelegateMessage()
	if err != nil {
		return HandleError(err)
	}

	// convert to model
	doc, ok := obj.(interfaces.Model)
	if !ok {
		return HandleError(errors.ErrorModelInvalidType)
	}

	// model delegate
	d := delegate.NewModelDelegate(doc)

	// apply method
	switch msg.GetMethod() {
	case interfaces.ModelDelegateMethodAdd:
		err = d.Add()
	case interfaces.ModelDelegateMethodSave:
		err = d.Save()
	case interfaces.ModelDelegateMethodDelete:
		err = d.Delete()
	case interfaces.ModelDelegateMethodGetArtifact, interfaces.ModelDelegateMethodRefresh:
		err = d.Refresh()
	}
	if err != nil {
		return HandleError(err)
	}

	// model
	m := d.GetModel()
	if msg.GetMethod() == interfaces.ModelDelegateMethodGetArtifact {
		m, err = d.GetArtifact()
		if err != nil {
			return nil, err
		}
	}

	// json bytes
	data, err := d.ToBytes(m)
	if err != nil {
		return nil, err
	}

	return HandleSuccessWithData(data)
}

func NewModelDelegateServer() (svr *ModelDelegateServer) {
	return &ModelDelegateServer{}
}
