package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/bogdan-user/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	l log.Logger
	greetpb.UnimplementedGreetServiceServer
}

func NewServer(log log.Logger) *Server {
	return &Server{l: log}
}

func (s *Server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("Greet function was inoked with %v\n", req)
	firstName := req.GetGreeting().FirstName
	result := "Hello " + firstName
	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

func (s *Server) ShowDetails(ctx context.Context, personRequest *greetpb.PersonDetailRequest) (*greetpb.PersonDetailResponse, error) {
	firstName := personRequest.GetPersonName()
	hobbies := []string{"Coding", "Working out"}

	fmt.Printf("Details for person: %v\n", firstName)
	return &greetpb.PersonDetailResponse{
		Employed: true,
		Age:      25,
		Hobbies:  hobbies,
	}, nil
}

func (s *Server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("GreetManyTimes invoked.")
	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)

		res := greetpb.GreetManyTimesResponse{
			Result: result,
		}

		stream.Send(&res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*Server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet invoked.")

	result := "Hello "
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.Send(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("error while receiving stream: %v\n", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
}

func (*Server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("BiDirectional invoked...")
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error while receiving the stream: %v\n", err)
			return err
		}

		firstName := req.Greeting.GetFirstName()
		lastName := req.Greeting.GetLastName()
		result := "Hello " + firstName + " " + lastName
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if err != nil {
			log.Fatalf("error while sending stream to client: %v\n", err)
			return err
		}

	}

}

func (*Server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {

	for i := 0; i < 6; i++ {
		if ctx.Err() == context.Canceled {
			//the client canceled the request
			fmt.Println("The client canceled the request")
			return nil, status.Error(codes.Canceled, "the client cancelled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.Greeting.GetFirstName()
	lastName := req.Greeting.GetLastName()
	result := "Hello " + firstName + " " + lastName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}

	return res, nil

}
