package main

import (
	"context"
	"log"
	"net"

	"github.com/gowo9/g3server"
	helloPB "github.com/gowo9/g3server/example/proto/hello"

	"google.golang.org/grpc"
)

type HelloController struct {
	helloPB.UnimplementedHelloServiceServer
}

func (hc *HelloController) Hello(context.Context, *helloPB.HelloRequest) (*helloPB.HelloResponse, error) {
	return &helloPB.HelloResponse{
		Msg: "hello",
	}, nil
}

func main() {
	// 建立一個 gRPC.Server
	gs := grpc.NewServer()
	helloPB.RegisterHelloServiceServer(gs, &HelloController{})

	// 將創建好的 gRPC.Server 和 gRPC-Gateway 的註冊函數作為參數傳給 ggserver.New
	s, err := g3server.New(context.Background(), gs, helloPB.RegisterHelloServiceHandler)
	if err != nil {
		log.Fatal("ggserver.New faile, ", err)
	}

	addr := "127.0.0.1:8080"
	tcpLis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("net.Listen faile, ", err)
	}

	log.Printf("start listen at %s\n", addr)
	if err = s.Serve(tcpLis); err != nil {
		log.Fatal("s.Serve faile, ", err)
	}
}
