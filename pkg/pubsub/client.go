package pubsub

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type Client struct {
	connection net.Conn
	pubsub *PubSub
	removeClient func(*Client)
	rw *bufio.ReadWriter
	lock *sync.Mutex //todo, benchmarking, unit tests
}

func newClient(connection net.Conn, pubsub *PubSub, removeClient func(*Client)) *Client {
	return &Client{
		connection,
		pubsub,
		removeClient,
		bufio.NewReadWriter(bufio.NewReader(connection), bufio.NewWriter(connection)),
		&sync.Mutex{},
	}
}

func (c *Client) Run() {
	c.receive()
}

func (c *Client) receive() {
	for {
		if bytes, _, err := c.rw.ReadLine(); err != nil {
			return
		} else {
			line := string(bytes)
			c.handleCommand(line)
		}
	}
}

func (c *Client) handleCommand(line string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	parts := strings.SplitN(line, " ", 3)

	if len(parts) == 0 {
		return
	}

	command := parts[0]

	switch command {
	case "pub":
		c.handlePublishCommand(parts)
		return

	case "sub":
		c.handleSubscribeCommand(parts)
		return

	case "kill":
		c.handleKillCommand()
		return
	}

	c.send("unknown command '%s'", command)
}

func (c *Client) handlePublishCommand(parts []string) {
	if len(parts) < 3 {
		c.send("expected: pub [topic] [message]")
		return
	}

	topic := parts[1]
	message := parts[2]
	if err := c.pubsub.Publish(topic, message); err != nil {
		c.send(err.Error())
	}
	c.send("published to topic [%s]", topic)
}

func (c *Client) handleSubscribeCommand(parts []string) {
	if len(parts) < 2 {
		c.send("expected: sub [topic]")
		return
	}

	topic := parts[1]
	subscription := c.pubsub.Subscribe(topic)
	go c.listenForMessages(subscription)
	c.send("subscribed to topic [%s]", topic)
}

func (c *Client) handleKillCommand() {
	os.Exit(0)
}

func (c *Client) listenForMessages(subscription *Subscription) {
	for {
		if message, ok := <- subscription.Channel(); !ok {
			return
		} else {
			c.lock.Lock()
			c.send("[%s] %s", subscription.Topic(), message)
			c.lock.Unlock()
		}
	}
}

func (c *Client) send(format string, args ...interface{}) {
	_, _ = c.rw.WriteString(fmt.Sprintln(fmt.Sprintf(format, args...)))
	_ = c.rw.Flush()
}

func (c *Client) Close() error {
	c.removeClient(c)
	return c.connection.Close()
}
