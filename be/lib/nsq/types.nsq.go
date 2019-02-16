package nsq

import nsq "github.com/nsqio/go-nsq"

// A Producer offers a standard interface for NSQ Producer functions
type Producer interface {
	Publish(string, []byte) error
}

// Consumer wrap consumer to interface to be testable
type Consumer interface {
	GetConsumer(topic, channel string, config *nsq.Config) (err error)
	ConnectNSQLookupD() error
	ConnectNSQD() error
	AddHandler(handler nsq.Handler, concurrency int)
}

// ConsumerTopicChannel is contain topic and channel
type ConsumerTopicChannel struct {
	Topic       string
	Channel     string
	HandlerFunc func(message *nsq.Message) error
	Concurrency int
}
