package grpc

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	pb "github.com/warthog618/modem/pkg/pb"
	"google.golang.org/grpc"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	clients    map[string]chan string
	mu         sync.Mutex
	authStream map[string]pb.AuthService_AuthKeyStreamServer
}

func NewAuthServiceServer() *AuthServiceServer {
	return &AuthServiceServer{
		clients:    make(map[string]chan string),
		authStream: make(map[string]pb.AuthService_AuthKeyStreamServer),
	}
}

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.clients[req.PhoneNumber]; !exists {
		s.clients[req.PhoneNumber] = make(chan string)
		return &pb.RegisterResponse{Success: true, Message: "Registered successfully"}, nil
	}

	return &pb.RegisterResponse{Success: false, Message: "Already registered"}, nil
}

func (s *AuthServiceServer) SendAuthKey(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ch, exists := s.clients[req.PhoneNumber]; exists {
		log.Printf("Sending auth key to %s", req.PhoneNumber)
		ch <- req.AuthKey
		return &pb.AuthResponse{Success: true, Message: "Auth key sent successfully"}, nil
	}

	return &pb.AuthResponse{Success: false, Message: "Phone number not registered"}, nil
}

func (s *AuthServiceServer) AuthKeyStream(stream pb.AuthService_AuthKeyStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving stream: %v", err)
			return err
		}

		s.mu.Lock()
		s.authStream[req.PhoneNumber] = stream
		s.mu.Unlock()

		go func(phoneNumber string) {
			for authKey := range s.clients[phoneNumber] {
				err := stream.Send(&pb.AuthKey{
					PhoneNumber: phoneNumber,
					AuthKey:     authKey,
				})
				if err != nil {
					log.Printf("Error sending auth key: %v", err)
					return
				}
			}
		}(req.PhoneNumber)
	}
}

func generateAuthKey() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	key := make([]byte, 10)
	for i := range key {
		key[i] = charset[rand.Intn(len(charset))]
	}
	return string(key)
}

func StartGRPCServer() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	authServiceServer := NewAuthServiceServer()
	pb.RegisterAuthServiceServer(s, authServiceServer)

	log.Printf("Server is listening on port :50051")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
