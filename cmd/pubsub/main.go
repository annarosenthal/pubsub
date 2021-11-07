package main

import "pubsub/pkg/pubsub"

func main() {
	server, err := pubsub.NewServer(8081)
	if err != nil {
		panic(err)
	}
	server.Listen()
}
