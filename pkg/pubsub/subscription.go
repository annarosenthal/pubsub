package pubsub

type Subscription struct {
	topic       string
	channel     chan string
	unsubscribe func(*Subscription)
}

func newSubscription(topic string, channel chan string, unsubscribe func(subscription *Subscription)) *Subscription {
	return &Subscription{topic, channel, unsubscribe}
}

func (s *Subscription) Close() error {
	s.unsubscribe(s)
	close(s.channel)
	return nil
}

func (s *Subscription) Channel() <-chan string {
	return s.channel
}

func (s *Subscription) Topic() string {
	return s.topic
}
