package pubsub

import (
	"fmt"
	"net"
)

type ClientConnectionSet map[*ClientConnection]interface{}

type Server struct {
	listener net.Listener
	pubsub   *PubSub
	clients  ClientConnectionSet
}

func NewServer() *Server {
	return &Server{
		pubsub:  NewPubSub(),
		clients: make(ClientConnectionSet),
	}
}

// Listen listens for new connections and processes commands from them, blocking until the Server is closed
func (s *Server) Listen(port int) error {
	if s.listener != nil {
		return fmt.Errorf("already listening")
	}

	if listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return err
	} else {
		fmt.Printf("Listening on port %d...\n", port)
		s.listener = listener
		s.acceptConnections() // block until Server is closed
		return nil
	}
}

// Close shuts down listener
func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) acceptConnections() {
	for {
		if conn, err := s.listener.Accept(); err != nil {
			return
		} else {
			fmt.Printf("client connected %s\n", conn.RemoteAddr())
			client := s.newClient(conn)
			go client.Run()
		}
	}
}

func (s *Server) newClient(conn net.Conn) *ClientConnection {
	client := newClient(conn, s.pubsub, s.clients.untrack)
	s.clients.track(client)
	return client
}

func (c ClientConnectionSet) track(client *ClientConnection) {
	c[client] = nil
}

func (c ClientConnectionSet) untrack(client *ClientConnection) {
	delete(c, client)
}
