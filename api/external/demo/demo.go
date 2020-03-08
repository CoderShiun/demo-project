package demo

import (
	"context"
	"demo_project/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

// DemoAPI defines the DemoAPI structure
type DemoAPI struct {
}

// DemoAPI returns the DemoAPI
func NewDemoAPI() *DemoAPI {
	return &DemoAPI{}
}

func (d *DemoAPI) Create(context.Context, *proto.CreateRequest) (*proto.CreateResponse, error) {
	return &proto.CreateResponse{
		Status: true,
	}, nil
}

func (d *DemoAPI) Get(context.Context, *proto.GetRequest) (*proto.GetResponse, error) {
	return &proto.GetResponse{
		Status: true,
	}, nil
}

func (d *DemoAPI) Delete(context.Context, *proto.DeleteRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil

}

func (d *DemoAPI) Update(context.Context, *proto.UpdateRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
