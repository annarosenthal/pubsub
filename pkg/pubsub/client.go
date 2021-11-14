package pubsub

import (
	"net"
	"pubsub/gen/pubsub"
)

type Client struct {
	connection net.Conn
	pr         *ProtoReader
	pw         *ProtoWriter
	messages   chan *Message
}

func NewClient(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	messages := make(chan *Message)
	client := &Client{conn, NewProtoReader(conn), NewProtoWriter(conn), messages}

	return client, nil
}

func (c *Client) Publish(topic string, message string) error {
	return c.pw.Write(&pubsub.Command{Topic: topic, Message: message, Action: pubsub.Command_PUB})
}

func (c *Client) Subscribe(topic string) error {
	return c.pw.Write(&pubsub.Command{Topic: topic, Action: pubsub.Command_SUB})
}

func (c *Client) Messages() <-chan *Message {
	return c.messages
}

func (c *Client) ProcessMessages() {
	for {
		var msg pubsub.Message
		if err := c.pr.Read(&msg); err != nil {
			return
		} else {
			c.handleMessage(&msg)
		}
	}
}

func (c *Client) Close() error {
	close(c.messages)
	return c.connection.Close()
}

func (c *Client) handleMessage(message *pubsub.Message) {

	c.messages <- &Message{message.Topic, message.Message}
}
