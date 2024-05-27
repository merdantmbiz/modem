package main

import (
	"log"
	"net"

	pb "github.com/warthog618/modem/gen"
	"github.com/warthog618/modem/pkg/config"
	rpc "github.com/warthog618/modem/pkg/grpc"
	"google.golang.org/grpc"
)

func main() {

	err := config.InitTomlConf("config", "./pkg/config")

	if err != nil {
		log.Println(err)
	}

	//start grpc stream
	lis, err := net.Listen("tcp", config.TomlConf.GRPC.PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs := rpc.NewAuthStreamService()

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, gs)
	pb.RegisterOTPServiceServer(s, gs)

	log.Printf("Server is running on port :%s", config.TomlConf.GRPC.PORT)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}