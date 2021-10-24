package pubsub

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
)

type ClientConnection struct {
	connection    net.Conn
	pubsub        *PubSub
	onClose       func(*ClientConnection)
	rw            *bufio.ReadWriter
	subscriptions map[*Subscription]interface{}
	channel 	  chan string
}

func newClient(connection net.Conn, pubsub *PubSub, onClose func(*ClientConnection)) *ClientConnection {
	return &ClientConnection{
		connection,
		pubsub,
		onClose,
		bufio.NewReadWriter(bufio.NewReader(connection), bufio.NewWriter(connection)),
		make(map[*Subscription]interface{}),
		make(chan string),
	}
}

func (c *ClientConnection) Run() error {
	go c.sender()
	for {
		if bytes, _, err := c.rw.ReadLine(); err != nil {
			return c.Close()
		} else {
			line := string(bytes)
			c.handleCommand(line)
		}
	}
}

func (c *ClientConnection) handleCommand(line string) {
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

	case "stats":
		c.handleStatsCommand()
		return

	case "quit":
		c.handleQuitCommand()
		return
	}
	c.send("unknown command '%s'", command)
}

func (c *ClientConnection) handlePublishCommand(parts []string) {
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

func (c *ClientConnection) handleSubscribeCommand(parts []string) {
	if len(parts) < 2 {
		c.send("expected: sub [topic]")
		return
	}

	topic := parts[1]
	subscription := c.pubsub.Subscribe(topic)
	c.subscriptions[subscription] = nil
	go c.listenForMessages(subscription)
	c.send("subscribed to topic [%s]", topic)
}

func (c *ClientConnection) handleKillCommand() {
	os.Exit(0)
}

func (c *ClientConnection) handleStatsCommand() {
	c.send("num go routines=%d", runtime.NumGoroutine())
}

func (c *ClientConnection) handleQuitCommand() {
	c.connection.Close()
}

func (c *ClientConnection) listenForMessages(subscription *Subscription) {
	for {
		if message, ok := <-subscription.Channel(); !ok {
			return
		} else {
			c.send("[%s] %s", subscription.Topic(), message)
		}
	}
}

func (c *ClientConnection) send(format string, args ...interface{}) {
	c.channel <- fmt.Sprintln(fmt.Sprintf(format, args...))
}

func (c *ClientConnection) sender() {
	for {
		if message, ok := <- c.channel; !ok {
			return
		} else {
			_, _ = c.rw.WriteString(message)
			_ = c.rw.Flush()
		}
	}
}

func (c *ClientConnection) Close() error {
	c.onClose(c)
	for subscription := range c.subscriptions {
		subscription.Close()
	}
	close(c.channel)
	return c.connection.Close()
}
