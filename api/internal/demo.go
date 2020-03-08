package internal

import (
	"context"
	"demo_project/proto"
)

// ApplicationAPI exports the grpc related functions.
type GrpcAPI struct {
}

// NewGrpcAPI creates a new grpc API.
func NewGrpcAPI() *GrpcAPI {
	return &GrpcAPI{
	}
}

func (g GrpcAPI) Get(context.Context, *proto.GetGrpcRequest) (*proto.GetGrpcResponse, error) {
	return &proto.GetGrpcResponse{Message:"Demo return message"}, nil
}