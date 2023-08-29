package main

import (
	"context"
	"log"
	"net"
	"net/http"
	pb "todo/proto/todo"
	service "todo/internal/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverPort = ":8080"
	gatewayPort = ":8090"
)

func startServer() {
	lis, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer()
	pb.RegisterTodoServer(s, &service.TodoServer{})
	log.Println("Serving gRPC on http://0.0.0.0" + serverPort)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()
}

func startGateway() {
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0" + serverPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterTodoHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    gatewayPort,
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0" + gatewayPort)
	log.Fatalln(gwServer.ListenAndServe())
}

func main() {
	startServer()
	startGateway()
}