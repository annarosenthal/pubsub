package main

import "pubsub/pkg/pubsub"

func main() {
	server := pubsub.NewServer()
	if err := server.Listen(8081); err != nil {
		panic(err)
	}
}
