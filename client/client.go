package main

import (
	"context"
	"fmt"
	"log"

	cntcharpb "github.com/nicewook/grpc-deadline/proto"
	"google.golang.org/grpc"
)

const (
	in1 = "abababab"
	in2 = "abcabcabcabc"
	in3 = "abbcccddddeeeee"
)

func main() {
	fmt.Println("Count Char gRPC client starts!")
	cc, dialErr := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if dialErr != nil {
		log.Fatalf("fail to Dial to gRPC server: %v", dialErr)
	}

	c := cntcharpb.NewCntCharServiceClient(cc)
	res, err := c.CntChar(context.Background(), &cntcharpb.CntCharReq{StrInput: in1})
	if err != nil {
		log.Printf("err from the gRPC server: %v\n", err)
	}
	fmt.Println("gRPC client got RPC res")
	fmt.Println(res.CntResult)

}
