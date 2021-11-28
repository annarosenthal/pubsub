package pubsub

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"pubsub/gen/pubsub"
)

type Server struct {
	pubsub.UnimplementedPubSubServer
	listener net.Listener
	pubsub   *PubSub
}

func NewServer(port int) (*Server, error) {
	if listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return nil, err
	} else {
		return &Server{
			pubsub.UnimplementedPubSubServer{},
			listener,
			NewPubSub(),
		}, nil
	}
}

// Listen listens for new connections and processes commands from them
func (s *Server) Listen() {
	fmt.Printf("Listening on %s...\n", s.listener.Addr().String())
	server := grpc.NewServer()
	pubsub.RegisterPubSubServer(server, s)
	server.Serve(s.listener)
}

// Close shuts down listener
func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s * Server) Subscribe(request *pubsub.SubscribeRequest, response pubsub.PubSub_SubscribeServer) error {
	subscription := s.pubsub.Subscribe(request.Topic)
	go func() {
		for {
			text, ok := <-subscription.Channel()
			if !ok {
				return
			}
			message := &pubsub.Message{Message: text}
			err := response.Send(message)
			if err != nil {
				return
			}
		}
	}()
	return nil
}

func (s * Server) Publish(_ context.Context, request *pubsub.PublishRequest) (*pubsub.PublishResponse, error) {
	err := s.pubsub.Publish(request.Topic, request.Message)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &pubsub.PublishResponse{}, nil
}
