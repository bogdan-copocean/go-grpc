package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/bogdan-user/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) Sum(ctx context.Context, sumRq *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {

	firstNum := sumRq.GetFirstNumber()
	secondNum := sumRq.GetSecondNumber()

	result := calculatorpb.SumResponse{
		SumResult: firstNum + secondNum,
	}

	return &result, nil
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading the stream: %v\n", err)
			return err
		}
		max := findMinAndMax(data.GetNumbers())

		err = stream.Send(&calculatorpb.FindMaximumResponse{
			Maximum: max,
		})
		if err != nil {
			log.Fatalf("error while sending data: %v\n", err)
			return err
		}
		fmt.Printf("Data sent: %v\n", max)
	}
	return nil

}

func (*server) SquareRoot(ctx context.Context, sqrequest *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := sqrequest.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received a negative number: %v\n", number))
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func findMinAndMax(a []int32) int32 {
	min := a[0]
	max := a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return max
}
