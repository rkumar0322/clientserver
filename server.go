package main

import (
	"clientserver/chat"
	// call the chat service we defined
	"fmt"
	"google.golang.org/grpc" // import the grpc package
	"log"
	"net"
)

func main() {
	fmt.Println("Our first gRPC server ")
	listener, err := net.Listen("tcp", "10.128.183.117:10000")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a new grpc server
	grpcServer := grpc.NewServer()

	ch := chat.Server{}
	chat.RegisterChatServiceServer(grpcServer, &ch)

	// register the endpoints we want to expose before serving this
	// over the existing TCP connection defined above
	err = grpcServer.Serve(listener)

	//display error in case of an error
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}