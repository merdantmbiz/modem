package grpc

import (
	"log"
	"sync"

	pb "github.com/warthog618/modem/gen" // Adjust the import path

	"google.golang.org/grpc/peer"
)

type tokenService interface {
	SendTokenToClient(clientID, token string) error
}
type Server struct {
	pb.UnimplementedExampleServiceServer
	mu            sync.Mutex
	Clients       map[string]pb.ExampleService_StreamDataServer
	tokenNotifier tokenService
}

func NewServer(tokenNotifier tokenService) *Server {
	return &Server{
		Clients:       make(map[string]pb.ExampleService_StreamDataServer),
		tokenNotifier: tokenNotifier,
	}
}
func (s *Server) StreamData(stream pb.ExampleService_StreamDataServer) error {
	p, _ := peer.FromContext(stream.Context())
	clientID := p.Addr.String()

	s.mu.Lock()
	s.Clients[clientID] = stream
	s.mu.Unlock()

	for {
		req, err := stream.Recv()
		if err != nil {
			s.mu.Lock()
			delete(s.Clients, clientID)
			s.mu.Unlock()
			return err
		}

		log.Printf("Received request from client ID: %s", req.ClientId)
	}
}

// SendTokenToClient sends a token to a specific client
func (s *Server) SendTokenToClient(clientID, token string) {
	s.mu.Lock()
	stream, ok := s.Clients[clientID]
	s.mu.Unlock()
	if !ok {
		log.Printf("Client not found: %s", clientID)
		return
	}

	response := &pb.StreamResponse{
		ClientId:     clientID,
		ResponseData: "Token: " + token,
	}
	if err := stream.Send(response); err != nil {
		log.Printf("Error sending token to client ID: %s, error: %v", clientID, err)
	} else {
		log.Printf("Token sent to client ID: %s", clientID)
	}
}
