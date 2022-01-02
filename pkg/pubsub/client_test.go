package pubsub

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pubsub/gen/pubsub"
	"pubsub/pkg/metrics"
	"testing"
	"time"
)

var (
	server *Server
	client *Client
)

func TestMain(m *testing.M) {
	var err error
	server, err = NewServer(8081, WithCollector(&metrics.StdOutCollector{}))
	if err != nil {
		panic(err)
	}
	go func() {
		defer server.Close()
		server.Listen()
	}()

	client, err = NewClient(":8081")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	m.Run()
}

func TestClient_Subscribe_Publish(t *testing.T) {
	err := client.Subscribe("topic")
	assert.NoError(t, err)

	err = client.Publish("topic", "message")
	assert.NoError(t, err)

	assertMessage(t, Topic, TestMessage)
}

func TestClient_Subscribe_Publish_WrongTopic(t *testing.T) {
	err := client.Subscribe("topic")
	assert.NoError(t, err)

	err = client.Publish("wrong", "message")
	assert.Error(t, err)
	assertStatus(t, err, codes.InvalidArgument)

	select {
	case msg := <-client.Messages():
		assert.Fail(t, "didn't expect a message")
		fmt.Printf("%+v", msg)
	case <-time.After(100 * time.Millisecond): // how can we not make it flakey
	}
}

func readMessage() *pubsub.Message {
	select {
	case msg := <-client.Messages():
		return msg
	case <-time.After(2 * time.Second):
		panic("timed out waiting for a message")
	}
}

func assertMessage(t *testing.T, topic string, message string) {
	msg := readMessage()
	assert.Equal(t, topic, msg.Topic)
	assert.Equal(t, message, msg.Message)
}

func assertStatus(t *testing.T, err error, code codes.Code) {
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, code, st.Code())
}
