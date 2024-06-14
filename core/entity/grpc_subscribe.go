package entity

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
)

type GrpcSubscribe struct {
	Stream   interfaces.GrpcStream
	Finished chan bool
}

func (sub *GrpcSubscribe) GetStream() interfaces.GrpcStream {
	return sub.Stream
}

func (sub *GrpcSubscribe) GetStreamBidirectional() interfaces.GrpcStreamBidirectional {
	stream, ok := sub.Stream.(interfaces.GrpcStreamBidirectional)
	if !ok {
		return nil
	}
	return stream
}

func (sub *GrpcSubscribe) GetFinished() chan bool {
	return sub.Finished
}
