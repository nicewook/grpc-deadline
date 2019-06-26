package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cntcharpb "github.com/nicewook/grpc-deadline/proto"
	"google.golang.org/grpc"
)

const in = "abbcccddddeeeee"

func main() {
	fmt.Println("Count Char gRPC client starts!")
	cc, dialErr := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if dialErr != nil {
		log.Fatalf("fail to Dial to gRPC server: %v", dialErr)
	}

	c := cntcharpb.NewCntCharServiceClient(cc)
	reqWithDeadline(c, 500*time.Millisecond, in)
	fmt.Println("--")
	time.Sleep(time.Second)

	reqWithDeadline(c, 100*time.Millisecond, in)
	fmt.Println("--")
	time.Sleep(time.Second)

	reqWithDeadlineCancel(c, 500*time.Millisecond, in)
	fmt.Println("--")

}

func reqWithDeadline(c cntcharpb.CntCharServiceClient, ms time.Duration, strInput string) {
	ctx, cancel := context.WithTimeout(context.Background(), ms)
	defer cancel()

	res, err := c.CntChar(ctx, &cntcharpb.CntCharReq{StrInput: strInput})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			code := statusErr.Code()
			rpcErr := statusErr.Err()
			if code == codes.DeadlineExceeded {
				fmt.Println(rpcErr)
			} else if code == codes.Canceled {
				fmt.Println(rpcErr)
			} else {
				fmt.Printf("unexpected gRPC error: %v\n", rpcErr)
			}
		} else {
			log.Fatalf("error while calling gRPC: %v", err)
		}
		return
	}
	fmt.Println("gRPC client got RPC res")
	fmt.Printf("%v\n", res.CntResult)
}

func reqWithDeadlineCancel(c cntcharpb.CntCharServiceClient, ms time.Duration, strInput string) {
	ctx, cancel := context.WithTimeout(context.Background(), ms)
	defer cancel()

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	res, err := c.CntChar(ctx, &cntcharpb.CntCharReq{StrInput: strInput})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			code := statusErr.Code()
			rpcErr := statusErr.Err()
			if code == codes.DeadlineExceeded {
				fmt.Println(rpcErr)
			} else if code == codes.Canceled {
				fmt.Println(rpcErr)
			} else {
				fmt.Printf("unexpected gRPC error: %v\n", rpcErr)
			}
		} else {
			log.Fatalf("error while calling gRPC: %v", err)
		}
		return
	}
	fmt.Println("gRPC client got RPC res")
	fmt.Println(res.CntResult)
}
