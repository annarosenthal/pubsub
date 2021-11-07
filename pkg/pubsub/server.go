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

func NewServer(port int) (*Server, error) {
	if listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return nil, err
	} else {
		return &Server{
			listener,
			NewPubSub(),
			make(ClientConnectionSet),
		}, nil
	}
}


// Listen listens for new connections and processes commands from them
func (s *Server) Listen() {
	fmt.Printf("Listening on %s...\n", s.listener.Addr().String())
	s.acceptConnections()
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
