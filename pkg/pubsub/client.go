package pubsub

import (
	"context"
	"google.golang.org/grpc"
	"pubsub/gen/pubsub"
)

type Client struct {
	connection grpc.ClientConnInterface
	messages   chan *Message
	client pubsub.PubSubClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	messages := make(chan *Message)
	client := pubsub.NewPubSubClient(conn)

	return &Client{conn, messages, client}, nil
}

func (c *Client) Publish(topic string, message string) error {
	request := &pubsub.PublishRequest{Topic: topic, Message: message}
	_, err := c.client.Publish(context.Background(),request)
	return err
}

func (c *Client) Subscribe(topic string) error {
	request := &pubsub.SubscribeRequest{Topic: topic}
	response, err := c.client.Subscribe(context.Background(),request)
	if err != nil {
		return err
	}

	go func() {
		for {
			message, err := response.Recv()
			if err != nil {
				return
			}
			c.messages <- &Message{Topic: topic, Text: message.Message}
		}
	}()
	return nil
}

func (c *Client) Messages() <-chan *Message {
	return c.messages
}

func (c *Client) Close() error {
	close(c.messages)
	return nil
}
