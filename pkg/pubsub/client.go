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

func (c *Client) ReceiveMessages() error {
	response, err := c.client.ReceiveMessages(context.Background(), &pubsub.ReceiveMessagesRequest{})
	if err != nil {
		return err
	}
	go func() {
		for {
			message, err := response.Recv()
			if err != nil {
				return
			}
			c.messages <- message
		}
	}()
	return nil
}

func (c *Client) Subscribe(topic string) error {
	request := &pubsub.SubscribeRequest{Topic: topic}
	_, err := c.client.Subscribe(context.Background(),request)
	return err
}

func (c *Client) Messages() <-chan *pubsub.Message {
	return c.messages
}

func (c *Client) Close() error {
	close(c.messages)
	return nil
}
