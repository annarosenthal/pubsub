package pubsub

import (
	"fmt"
	"net"
)

type Server struct {
	listener net.Listener
	pubsub   *PubSub
	clients map[*ClientConnection]interface{}
}

func NewServer() *Server {
	return &Server{
		pubsub: NewPubSub(),
		clients: make(map[*ClientConnection]interface{}),
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
		s.doWork() // block until Server is closed
		return nil
	}
}

func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) doWork() {
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
	client := newClient(conn, s.pubsub, s.removeClient)
	s.clients[client] = nil
	return client
}

func (s *Server) removeClient(client *ClientConnection) {
	delete(s.clients, client)
}
