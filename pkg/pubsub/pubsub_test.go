package pubsub

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const Topic = "topic"
const Message = "message"

func TestPubSub_PublishError(t *testing.T) {
	pubsub := NewPubSub()
	err := pubsub.Publish(Topic, Message)
	assert.NotNil(t, err)
}

func TestPubSub_Publish(t *testing.T) {
	pubsub := NewPubSub()
	subscription := pubsub.Subscribe(Topic)
	assert.NotNil(t, subscription)


	go func () {
		err := pubsub.Publish(Topic, Message)
		assert.Nil(t, err)
	} ()

	message, ok := <-subscription.Channel()
	assert.True(t, ok)
	assert.Equal(t, Message, message)
}

func TestPubSub_PublishCloseAndPublish(t *testing.B) {
	pubsub := NewPubSub()

	subscription := pubsub.Subscribe(Topic)
	assert.NotNil(t, subscription)
	_ = subscription.Close()

	go func () {
		err := pubsub.Publish(Topic, Message)
		assert.NotNil(t, err)
	} ()

	_, ok := <-subscription.Channel()
	assert.False(t, ok)
}
