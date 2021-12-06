package pubsub

import "pubsub/pkg/metrics"

type ServerOption func(*Server)

func WithCollector(collector metrics.Collector) ServerOption {
	return func(server *Server) {
		server.collector = collector
	}
}
