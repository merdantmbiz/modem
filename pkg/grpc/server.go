package grpc

import (
	"log"
	"sync"

	pb "github.com/warthog618/modem/gen" // Adjust the import path
)

type AuthStreamService interface {
	SendTokenToClient(clientID, token string)
}

type server struct {
	pb.UnimplementedAuthServiceServer
	mu      sync.Mutex
	Clients map[string]pb.AuthService_AuthStreamServer
}

func (s *server) StreamData(stream pb.AuthService_AuthStreamServer) error {
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

// SendTokenToClient sends a token to a specific client
func (s *server) SendTokenToClient(clientID, token string) {
	s.mu.Lock()
	stream, ok := s.Clients[clientID]
	s.mu.Unlock()
	if !ok {
		log.Println(s.Clients)
		log.Printf("Client not found: %s clients count: %d", clientID, len(s.Clients))
		return
	}

	response := &pb.AuthResponse{
		ClientId:     clientID,
		ResponseData: token,
	}
	if err := stream.Send(response); err != nil {
		log.Printf("Error sending token to client ID: %s, error: %v", clientID, err)
		s.mu.Lock()
		delete(s.Clients, clientID)
		s.mu.Unlock()
	} else {
		log.Printf("Token sent to client ID: %s", clientID)
	}
}
