package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	helloPB "github.com/gowo9/g3server/example/proto/hello"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serverAddr = "127.0.0.1:8080"

func runRESTRequest() {
	//nolint:noctx // just for a simple example
	rsp, err := http.Get(fmt.Sprintf("http://%s/", serverAddr))
	if err != nil {
		log.Println("http.Get failed,", err)
		return
	}
	defer rsp.Body.Close()

	rspBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		log.Println("io.ReadAll failed,", err)
		return
	}

	fmt.Printf("REST response: %s\n", string(rspBytes))
}

func runGRPCRequest() {
	conn, err := grpc.Dial(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	helloClient := helloPB.NewHelloServiceClient(conn)

	rsp, err := helloClient.Hello(context.Background(), &helloPB.HelloRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("gRPC response: %+v\n", rsp)
}

func main() {
	runRESTRequest()
	runGRPCRequest()
}
