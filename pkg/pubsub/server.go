package pubsub

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"pubsub/gen/pubsub"
	"pubsub/pkg/metrics"
)

type Server struct {
	pubsub.UnimplementedPubSubServer
	listener net.Listener
	pubsub   *PubSub
	collector metrics.Collector
}

func NewServer(port int, options ...ServerOption) (*Server, error) {
	if listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		return nil, err
	} else {
		server := &Server{
			pubsub.UnimplementedPubSubServer{},
			listener,
			NewPubSub(),
			&metrics.DefaultCollector{},
		}
		for _, option := range options {
			option(server)
		}
		return server, nil
	}
}

// Listen listens for new connections and processes commands from them
func (s *Server) Listen() error {
	fmt.Printf("Listening on %s...\n", s.listener.Addr().String())
	server := grpc.NewServer()
	pubsub.RegisterPubSubServer(server, s)
	return server.Serve(s.listener)
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
	s.collector.Increment(metrics.SubscribeCountMetric,tags(request.Topic))
	if err := response.Send(&pubsub.Message{}); err!= nil {
		s.collector.Increment(metrics.SubscribeErrorCountMetric,tags(request.Topic))
		return err
	}
	for {
		text, ok := <-subscription.Channel()
		if !ok {
			return nil
		}
		message := &pubsub.Message{Message: text, Topic: request.Topic}
		if err := response.Send(message); err != nil {
			if err != io.EOF {
				s.collector.Increment(metrics.SubscribeErrorCountMetric,tags(request.Topic))
			}
			return err
		}
	}
}

func (s * Server) Publish(_ context.Context, message *pubsub.Message) (*pubsub.PublishResponse, error) {
	err := s.pubsub.Publish(message.Topic, message.Message)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &pubsub.PublishResponse{}, nil
}

func tags(topic string) metrics.Tags {
	return metrics.Tags{
		"topic": topic,
	}
}
