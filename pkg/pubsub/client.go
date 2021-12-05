package pubsub

import (
	"context"
	"google.golang.org/grpc"
	"pubsub/gen/pubsub"
)

type Client struct {
	connection grpc.ClientConnInterface
	messages   chan *pubsub.Message
	client pubsub.PubSubClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	messages := make(chan *pubsub.Message)
	client := pubsub.NewPubSubClient(conn)

	return &Client{conn, messages, client}, nil
}

func (c *Client) Publish(topic string, message string) error {
	request := &pubsub.Message{Topic: topic, Message: message}
	_, err := c.client.Publish(context.Background(),request)
	return err
}

func (c *Client) Subscribe(topic string) error {
	request := &pubsub.SubscribeRequest{Topic: topic}
	response, err := c.client.Subscribe(context.Background(),request)
	if err != nil {
		return err
	}
	_, err = response.Recv()

	go func() {
		for {
			if msg, err := response.Recv(); err != nil {
				return
			} else {
				c.messages <- msg
			}

		}
	}()
	return err
}

func (c *Client) Messages() <-chan *pubsub.Message {
	return c.messages
}

func (c *Client) Close() error {
	close(c.messages)
	return nil
}
