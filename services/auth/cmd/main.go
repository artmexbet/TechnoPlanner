package main

import (
	"net"

	"auth/internal/server"
	"auth/internal/service"
	"proto"

	"google.golang.org/grpc"
)

func main() {
	svc := service.NewAuth()
	handler := server.NewHandler(svc)

	grpcServer := grpc.NewServer()
	proto.RegisterAuthServer(grpcServer, handler)
	grpcListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	if err = grpcServer.Serve(grpcListener); err != nil {
		panic(err) // TODO: handle error properly
	}
}
