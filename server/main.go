package main

import (
	"fmt"
	"go-grpc-restaurant-server/grpc_server"
	"go-grpc-restaurant-server/logger"
	"go-grpc-restaurant-server/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":8001"
)

func init() {
	log.SetPrefix("[LOG] ")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
}

func main() {
	listener, err := net.Listen("tcp", port)

	if err != nil {
		logger.Logger.Error(fmt.Errorf("failed to start grpc server: %v", err).Error())
	}

	grpcServer := grpc.NewServer()

	proto.RegisterOrderServiceServer(grpcServer, &grpc_server.OrderServiceServer{})
	logger.Logger.Info(fmt.Sprintf("server started at port %s", port))

	err = grpcServer.Serve(listener)

	if err != nil {
		logger.Logger.Fatal(fmt.Errorf("failed to start grpc server: %v", err).Error())
	}
}
