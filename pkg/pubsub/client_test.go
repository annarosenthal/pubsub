package pubsub

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	server *Server
	client *Client
)

func TestMain(m *testing.M) {
	var err error
	server, err = NewServer(12345)
	if err != nil {
		panic(err)
	}
	go func() {
		defer server.Close()
		server.Listen()
	}()

	client, err = NewClient(":12345")
	if err != nil {
		panic(err)
	}
	go func() {
		defer client.Close()
		client.ProcessMessages()
	}()

	m.Run()
}

func TestClient_Subscribe_Publish(t *testing.T) {
	err := client.Subscribe("topic")
	assert.NoError(t, err)
	err = client.Publish("topic", "message")
	assert.NoError(t, err)

	msg := readMessage()
	assert.Equal(t, "topic", msg.Topic)
	assert.Equal(t, "message", msg.Text)
}

func TestClient_Subscribe_Publish_WrongTopic(t *testing.T) {
	err := client.Subscribe("topic")
	assert.NoError(t, err)
	err = client.Publish("wrong", "message")
	assert.NoError(t, err)

	select {
	case <-client.Messages():
		assert.Fail(t, "didn't expect a message")
	case <-time.After(100 * time.Millisecond): // how can we not make it flakey
	}
}

func readMessage() *Message {
	select {
	case msg := <-client.Messages():
		return msg
	case <-time.After(2 * time.Second):
		panic("timed out waiting for a message")
	}
}
