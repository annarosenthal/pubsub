package pubsub

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	connection net.Conn
	rw *bufio.ReadWriter
	messages chan *Message
}

func NewClient (address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	rw := bufio.NewReadWriter(bufio.NewReader(conn),bufio.NewWriter(conn))
	messages := make(chan *Message)
	client := &Client{conn, rw, messages }

	return client, nil
}

func (c *Client) Publish(topic string, message string) error {
	_, err := c.rw.WriteString( fmt.Sprintf("pub %s %s\n", topic, message))
	if err != nil {
		return err
	}
	return c.rw.Flush()
}

func (c *Client) Subscribe(topic string) error {
	_, err := c.rw.WriteString( fmt.Sprintf("sub %s\n", topic))
	if err != nil {
		return err
	}
	return c.rw.Flush()
}

func (c *Client) Messages() <-chan *Message {
	return c.messages
}

func (c *Client) ProcessMessages()  {
	for {
		if bytes, _, err := c.rw.ReadLine(); err != nil {
			return
		} else {
			line := string(bytes)
			c.handleMessage(line)
		}
	}
}

func (c *Client) Close() error {
	close(c.messages)
	return c.connection.Close()
}

func (c *Client) handleMessage(line string) {
	parts := strings.SplitN(line, " ",2)
	if len(parts) < 2 {
		return
	}
	if !strings.HasPrefix(parts[0], "[") || !strings.HasSuffix(parts[0], "]") {
		return
	}
	topic := parts[0][1:len(parts[0])-1]
	message := parts[1]

	c.messages <- &Message{topic, message}
}
