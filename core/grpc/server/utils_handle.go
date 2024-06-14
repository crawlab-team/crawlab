package server

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
)

func HandleError(err error) (res *grpc.Response, err2 error) {
	trace.PrintError(err)
	return &grpc.Response{
		Code:  grpc.ResponseCode_ERROR,
		Error: err.Error(),
	}, err
}

func HandleSuccess() (res *grpc.Response, err error) {
	return &grpc.Response{
		Code:    grpc.ResponseCode_OK,
		Message: "success",
	}, nil
}

func HandleSuccessWithData(data interface{}) (res *grpc.Response, err error) {
	var bytes []byte
	switch data.(type) {
	case []byte:
		bytes = data.([]byte)
	default:
		bytes, err = json.Marshal(data)
		if err != nil {
			return HandleError(err)
		}
	}
	return &grpc.Response{
		Code:    grpc.ResponseCode_OK,
		Message: "success",
		Data:    bytes,
	}, nil
}

func HandleSuccessWithListData(data interface{}, total int) (res *grpc.Response, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return HandleError(err)
	}
	return &grpc.Response{
		Code:    grpc.ResponseCode_OK,
		Message: "success",
		Data:    bytes,
		Total:   int64(total),
	}, nil
}
