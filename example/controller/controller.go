package controller

import (
	"context"

	helloPB "github.com/gowo9/g3server/example/proto/hello"
)

type HelloController struct {
	helloPB.UnimplementedHelloServiceServer
}

func (hc *HelloController) Hello(context.Context, *helloPB.HelloRequest) (*helloPB.HelloResponse, error) {
	return &helloPB.HelloResponse{
		Msg: "hello",
	}, nil
}
