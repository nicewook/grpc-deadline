package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	cntcharpb "github.com/nicewook/grpc-deadline/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) CntChar(ctx context.Context, req *cntcharpb.CntCharReq) (*cntcharpb.CntCharRes, error) {
	fmt.Println("gRPC server received RPC req")

	cntCharMap := make(map[string]int)
	msg := req.GetStrInput()
	for _, char := range msg {
		cntCharMap[string(char)]++
	}

	canceled := make(chan string, 1)
	go func() {
		for {
			if ctx.Err() == context.Canceled {
				canceled <- "Client canceled"
			}
			time.Sleep(time.Millisecond)
		}
	}()

	select {
	case <-canceled:
		fmt.Print("client canceled the request\n--\n")
		return nil, status.Error(codes.Canceled, "client canceled the request")

	case <-time.After(300 * time.Millisecond):
		fmt.Print("server response the rpc call\n--\n")
		cntResult := fmt.Sprintf("Input: %s\nCount char result:\n%v", msg, cntCharMap)
		return &cntcharpb.CntCharRes{CntResult: cntResult}, nil
	}
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
