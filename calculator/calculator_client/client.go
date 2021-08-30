package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bogdan-user/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello, I'm a client")

	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to: %v\n", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)
	doErrUnary(c)
	// doBiDirectional(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	log.Println("Calculator client")

	calcRequest := calculatorpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 12,
	}

	res, err := c.Sum(context.Background(), &calcRequest)

	if err != nil {
		log.Fatalf("Sum RPC error: %v\n", err)
	}

	fmt.Println(res.SumResult)
}

func doErrUnary(c calculatorpb.CalculatorServiceClient) {
	log.Println("Error unary invoked")

	calcRequest := calculatorpb.SquareRootRequest{
		Number: 9,
	}

	res, err := c.SquareRoot(context.Background(), &calcRequest)

	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC
			fmt.Println()
			fmt.Printf("Error message from server: %v\n", resErr.Message())
			fmt.Printf("Code Error: %v\n", resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number")
				return
			}
		} else {
			log.Fatalf("Bigger error calling sqroot: %v\n", err)
			return
		}
	}

	fmt.Printf("Result of square root: %v\n", res.NumberRoot)
}

func doBiDirectional(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("BiDirectional started...")
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while creating the stream: %v\n", err)
	}
	waitc := make(chan struct{})

	requests := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Numbers: []int32{1, 2, 3, 4, 5},
		},
		&calculatorpb.FindMaximumRequest{
			Numbers: []int32{10, 222, 13, 4124, 35},
		},
		&calculatorpb.FindMaximumRequest{
			Numbers: []int32{1110, 22, 143, 124, 325},
		},
		&calculatorpb.FindMaximumRequest{
			Numbers: []int32{1033, 52, 13512, 4, 635},
		},
	}

	go func() {
		// func to send bunch of msgs
		for _, request := range requests {
			fmt.Printf("sending the message: %v\n", request)
			stream.Send(request)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		// 	func to receive bunch of msgs
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("error while receiving: %v\n", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetMaximum())
		}
		close(waitc)

	}()

	<-waitc
}
