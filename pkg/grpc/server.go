package grpc

import (
	"context"
	"fmt"
	"log"
	"sync"

	pb "github.com/warthog618/modem/gen" // Adjust the import path
	"github.com/warthog618/modem/pkg/jwt"
)

type IAuthService interface {
	pb.AuthServiceServer
	pb.OTPServiceServer
}

type authStreamServer struct {
	pb.UnimplementedAuthServiceServer
	pb.UnimplementedOTPServiceServer
	mu      sync.Mutex
	Clients map[string]pb.AuthService_AuthStreamServer
}

func NewAuthStreamService() IAuthService {
	return &authStreamServer{
		Clients: make(map[string]pb.AuthService_AuthStreamServer),
	}
}
func (s *authStreamServer) AuthStream(stream pb.AuthService_AuthStreamServer) error {
	// p, _ := peer.FromContext(stream.Context())
	// clientID := p.Addr.String()
	log.Printf("total clients before: %d", len(s.Clients))

	for {
		req, err := stream.Recv()
		if err != nil {
			// s.mu.Lock()
			// delete(s.Clients, clientID)
			// s.mu.Unlock()
			return err
		}
		s.mu.Lock()
		s.Clients[req.ClientId] = stream
		s.mu.Unlock()

		log.Printf("Received request from client ID: %s clients after: %d", req.ClientId, len(s.Clients))
	}
}

func (s *authStreamServer) PassOTP(ctx context.Context, req *pb.OTPRequest) (*pb.OTPResponse, error) {
	token, err := jwt.GenerateJWT(req.ClientId)
	s.mu.Lock()
	stream, ok := s.Clients[req.ClientId]
	s.mu.Unlock()

	if !ok {
		return &pb.OTPResponse{
			ClientId: fmt.Sprintf("Client not found: %s clients count: %d", req.ClientId, len(s.Clients)),
		}, err
	}

	response := &pb.AuthResponse{
		ClientId:     req.ClientId,
		ResponseData: token,
	}
	if err := stream.Send(response); err != nil {
		log.Printf("Error sending token to client ID: %s, error: %v", req.ClientId, err)
		s.mu.Lock()
		delete(s.Clients, req.ClientId)
		s.mu.Unlock()
	} else {
		log.Printf("Token sent to client ID: %s", req.ClientId)
	}

	return &pb.OTPResponse{
		ClientId: req.ClientId,
	}, err
}
