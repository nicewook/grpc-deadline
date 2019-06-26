package main

import (
	"context"
	"fmt"
	"log"
	"net"

	cntcharpb "github.com/nicewook/grpc-deadline/proto"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) CntChar(cnt context.Context, req *cntcharpb.CntCharReq) (*cntcharpb.CntCharRes, error) {
	fmt.Println("gRPC server received RPC req")
	return &cntcharpb.CntCharRes{}, nil
}

func main() {
	fmt.Println("Count Char gRPC server starts!")
	l, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("can not net.Listen: %v", err)
	}

	s := grpc.NewServer()
	cntcharpb.RegisterCntCharServiceServer(s, &server{})

	if serveErr := s.Serve(l); err != nil {
		log.Fatalf("can not serve gRPC server: %v", serveErr)
	}
}
