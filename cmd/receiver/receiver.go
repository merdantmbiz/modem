package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/warthog618/modem/gen"

	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewAuthServiceClient(conn)

	// Call the StreamData method
	stream, err := client.AuthStream(context.Background())
	if err != nil {
		log.Fatalf("could not open stream: %v", err)
	}

	waitc := make(chan struct{})

	// Send messages to the stream in a goroutine
	go func() {
		req := &pb.AuthRequest{
			ClientId:    "+99363432211",
			RequestData: fmt.Sprintf("request_data_%d", 1),
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("failed to send a request: %v", err)
		}
		time.Sleep(time.Minute)
		stream.CloseSend()
	}()

	// Receive messages from the stream
	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				log.Fatalf("failed to receive a response: %v", err)
				close(waitc)
				return
			}
			fmt.Printf("Received: %s\n", res.ResponseData)
		}
	}()

	<-waitc
}
