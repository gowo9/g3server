package main

import (
	"context"
	"log"
	"net"

	"github.com/gowo9/g3server"
	"github.com/gowo9/g3server/example/controller"
	helloPB "github.com/gowo9/g3server/example/proto/hello"

	"google.golang.org/grpc"
)

func main() {
	// 建立一個 gRPC.Server
	gs := grpc.NewServer()
	helloPB.RegisterHelloServiceServer(gs, &controller.HelloController{})

	// 將創建好的 gRPC.Server 和 gRPC-Gateway 的註冊函數作為參數傳給 ggserver.New
	s, err := g3server.New(context.Background(), gs, helloPB.RegisterHelloServiceHandler)
	if err != nil {
		log.Fatal("ggserver.New faile, ", err)
	}

	tcpLis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("net.Listen faile, ", err)
	}

	if err = s.Serve(tcpLis); err != nil {
		log.Fatal("s.Serve faile, ", err)
	}
}
