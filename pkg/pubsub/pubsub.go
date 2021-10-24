package pubsub

import "fmt"

type PubSub struct {
	subscriptions map[string]map[*Subscription]interface{} // topic to a set of subscriptions
}

func NewPubSub() *PubSub {
	return &PubSub{subscriptions: make(map[string]map[*Subscription]interface{})}
}

// Publish sends a message to a topic
func (p *PubSub) Publish(topic string, message string) error {
	if subscriptions, ok := p.subscriptions[topic]; ok {
		for subscription := range subscriptions {
			subscription.channel <- message
		}
		return nil
	} else {
		return fmt.Errorf("failed to publish to topic %s because there are no subscribers", topic)
	}
}

func (p *PubSub) Subscribe(topic string) *Subscription {
	var subscriptions map[*Subscription]interface{}
	var ok bool

	if subscriptions, ok = p.subscriptions[topic]; !ok {
		subscriptions = make(map[*Subscription]interface{})
	}

	subscription := newSubscription(topic, make(chan string), p.unsubscribe)

	subscriptions[subscription] = nil
	p.subscriptions[topic] = subscriptions

	return subscription
}

func (p *PubSub) unsubscribe(subscription *Subscription) {
	if subscriptions, ok := p.subscriptions[subscription.Topic()]; ok {
		delete(subscriptions, subscription)
	}
}
