package main

import (
	"pubsub/pkg/metrics"
	"pubsub/pkg/pubsub"
)

func main() {
	collector := &metrics.PrometheusCollector{}
	server, err := pubsub.NewServer(8081, pubsub.WithCollector(collector))
	if err != nil {
		panic(err)
	}
	server.Listen()
}
