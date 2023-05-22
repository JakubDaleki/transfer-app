package main

import (
	"fmt"
	"net"

	"github.com/JakubDaleki/transfer-app/query-service/db"
	"github.com/JakubDaleki/transfer-app/query-service/server"
	pb "github.com/JakubDaleki/transfer-app/shared-dependencies/grpc"
	"google.golang.org/grpc"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		fmt.Println(err)
	}
	grpcServer := grpc.NewServer()
	s := &server.GreeterService{Db: db}
	pb.RegisterGreeterServer(grpcServer, s)
	grpcServer.Serve(lis)
}
