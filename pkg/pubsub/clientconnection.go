package pubsub

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
)

type ClientConnection struct {
	connection net.Conn
	pubsub *PubSub
	removeClient func(*ClientConnection)
	rw *bufio.ReadWriter
	lock *sync.Mutex //todo, benchmarking, unit tests
	subscriptions map[*Subscription]interface{}
}

func newClient(connection net.Conn, pubsub *PubSub, removeClient func(*ClientConnection)) *ClientConnection {
	return &ClientConnection{
		connection,
		pubsub,
		removeClient,
		bufio.NewReadWriter(bufio.NewReader(connection), bufio.NewWriter(connection)),
		&sync.Mutex{},
		make(map[*Subscription]interface{}),
	}
}

func (c *ClientConnection) Run() {
	c.receive()
}

func (c *ClientConnection) receive() {
	for {
		if bytes, _, err := c.rw.ReadLine(); err != nil {
			c.Close()
		} else {
			line := string(bytes)
			c.handleCommand(line)
		}
	}
}

func (c *ClientConnection) handleCommand(line string) {
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
	c.Close()
}

func (c *ClientConnection) listenForMessages(subscription *Subscription) {
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

func (c *ClientConnection) send(format string, args ...interface{}) {
	_, _ = c.rw.WriteString(fmt.Sprintln(fmt.Sprintf(format, args...)))
	_ = c.rw.Flush()
}

func (c *ClientConnection) Close() error {
	c.removeClient(c)
	for subscription := range c.subscriptions {
		subscription.Close()
	}
	return c.connection.Close()
}
