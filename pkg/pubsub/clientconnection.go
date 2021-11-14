package pubsub

import (
	"fmt"
	"net"
	"pubsub/gen/pubsub"
)

type ClientConnection struct {
	connection    net.Conn
	pubsub        *PubSub
	onClose       func(*ClientConnection)
	pr            *ProtoReader
	pw            *ProtoWriter
	subscriptions map[*Subscription]interface{}
	channel       chan *pubsub.Message
}

func newClient(connection net.Conn, pb *PubSub, onClose func(*ClientConnection)) *ClientConnection {
	return &ClientConnection{
		connection,
		pb,
		onClose,
		NewProtoReader(connection),
		NewProtoWriter(connection),
		make(map[*Subscription]interface{}),
		make(chan *pubsub.Message),
	}
}

func (c *ClientConnection) Run() error {
	go c.sender()
	for {
		var cmd pubsub.Command
		if err := c.pr.Read(&cmd); err != nil {
			return c.Close()
		} else {
			c.handleCommand(cmd)
		}
	}
}

func (c *ClientConnection) handleCommand(cmd pubsub.Command) {

	switch cmd.Action {
	case pubsub.Command_PUB:
		c.handlePublishCommand(cmd.Topic, cmd.Message)
		return

	case pubsub.Command_SUB:
		c.handleSubscribeCommand(cmd.Topic)
		return
	default:
		c.error("unknown command")
	}

}

func (c *ClientConnection) handlePublishCommand(topic string, message string) {
	if err := c.pubsub.Publish(topic, message); err != nil {
		c.error(err.Error())
	} else {
		c.info("published to topic [%s]", topic)
	}
}

func (c *ClientConnection) handleSubscribeCommand(topic string) {
	subscription := c.pubsub.Subscribe(topic)
	c.subscriptions[subscription] = nil
	go c.listenForMessages(subscription)
	c.info("subscribed to topic [%s]", topic)
}

func (c *ClientConnection) listenForMessages(subscription *Subscription) {
	for {
		if message, ok := <-subscription.Channel(); !ok {
			return
		} else {
			c.channel <- &pubsub.Message{Topic: subscription.Topic(), Message: message}
		}
	}
}

func (c *ClientConnection) info(format string, args ...interface{}) {
	c.channel <- &pubsub.Message{Topic: "info", Message: fmt.Sprintf(format, args...)}
}

func (c *ClientConnection) error(format string, args ...interface{}) {
	c.channel <- &pubsub.Message{Topic: "error", Message: fmt.Sprintf(format, args...)}
}


func (c *ClientConnection) sender() {
	for {
		if message, ok := <-c.channel; !ok {
			return
		} else {
			_ = c.pw.Write(message)
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
