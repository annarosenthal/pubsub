package pubsub

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const Topic = "topic"
const TestMessage = "message"

func TestPubSub_PublishError(t *testing.T) {
	pubsub := NewPubSub()
	err := pubsub.Publish(Topic, TestMessage)
	assert.NotNil(t, err)
}

func TestPubSub_Publish(t *testing.T) {
	pubsub := NewPubSub()
	subscription := pubsub.Subscribe(Topic)
	assert.NotNil(t, subscription)

	go func() {
		err := pubsub.Publish(Topic, TestMessage)
		assert.Nil(t, err)
	}()

	message, ok := <-subscription.Channel()
	assert.True(t, ok)
	assert.Equal(t, TestMessage, message)
}

func TestPubSub_PublishCloseAndPublish(t *testing.T) {
	pubsub := NewPubSub()

	subscription := pubsub.Subscribe(Topic)
	assert.NotNil(t, subscription)
	_ = subscription.Close()

	go func() {
		err := pubsub.Publish(Topic, TestMessage)
		assert.NotNil(t, err)
	}()

	_, ok := <-subscription.Channel()
	assert.False(t, ok)
}

func benchmarkPubSub_Publish(num int, b *testing.B) {
	pubsub := NewPubSub()
	subscriptions := make([]*Subscription, num)

	for i := 0; i < num; i++ {
		subscription := pubsub.Subscribe(Topic)
		subscriptions[i] = subscription

		go func() {
			for {
				if _, ok := <-subscription.channel; !ok {
					return
				}
			}
		}()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = pubsub.Publish(Topic, TestMessage)
	}

	b.StopTimer()
	for _, subscription := range subscriptions {
		subscription.Close()
	}
}

func BenchmarkPubSub_Publish10(b *testing.B) {
	benchmarkPubSub_Publish(1000, b)
}
