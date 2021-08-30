package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bogdan-user/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello, I'm a client")

	tls := true
	opts := grpc.WithInsecure()

	if tls {
		certFile := "ssl/ca.crt" // Certificate Authority Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")

		if sslErr != nil {
			log.Fatalf("Error while loading SSL certificate: %v\n", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect to: %v\n", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)
	// doUnaryTwo(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDirectionalStreaming(c)
	// doUnaryWithDeadline(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do Unary RPC")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bogdan",
			LastName:  "Copocean",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v\n", err)
	}

	log.Printf("Response from Greet RPC: %v\n", res.Result)
}

func doUnaryTwo(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do Unary (TWO) RPC")

	user := greetpb.PersonDetailRequest{
		PersonName: "Bogdan-Lucian Copocean",
	}

	res, err := c.ShowDetails(context.Background(), &user)
	if err != nil {
		log.Fatalf("error while calling ShowDetails RPC: %v\n", err)
	}

	log.Printf("Response from ShowDetails RPC: %v, %v, %v\n", res.Hobbies, res.Age, res.Employed)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting server streaming RPC...")

	req := greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bogdan",
			LastName:  "Copocean",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v\n", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v\n", err)
		}
		log.Printf("Response from GreetManyTimes: %v\n", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting client streaming...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bogdan",
				LastName:  "Copocean",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Roxana",
				LastName:  "Copocean",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Roxana",
				LastName:  "Mica",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bogdan",
				LastName:  "Mica",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while long greet: %v\n", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	err = stream.CloseSend()
	if err != nil {
		log.Fatalf("error while closing the stream: %v\n", err)
	}

	res, err := stream.Recv()
	if err != nil {
		log.Fatalf("error while receiving stream: %v\n", err)
	}

	fmt.Printf("Result: %v\n", res.Result)
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("BiDirectional started...")
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating the stream: %v\n", err)
	}

	waitc := make(chan struct{})
	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bogdan",
				LastName:  "Copocean",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Roxana",
				LastName:  "Copocean",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Roxana",
				LastName:  "Mica",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bogdan",
				LastName:  "Mica",
			},
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
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)

	}()

	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to doUnaryWithDeadline RPC...")

	req := greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bogdan",
			LastName:  "Lucian",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, &req)

	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout hit! Deadline has exceeded")
			} else {
				fmt.Printf("Unexpected error: %v\n", statusErr)
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadline RPC: %v\n", err)
		}

		log.Fatalf("error while calling Greetwithdeadline RPC: %v\n", err)
		return
	}

	log.Printf("Response from server: %v\n", res.Result)

}
