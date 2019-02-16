package gonsq

import (
	"log"

	gnsq "github.com/nsqio/go-nsq"
	"github.com/template/be/lib/nsq"
)

var (
	newProducer = gnsq.NewProducer
)

// NewProducer GoNSQ module
func NewProducer(config Config) (nsq.Producer, error) {

	g := &gonsq{}

	// Create new producer
	var err error
	if g.prod, err = newProducer(config.Endpoint, gnsq.NewConfig()); err != nil {
		return nil, err
	}
	return g, nil
}

// NewConsumer get instance of consumer
func NewConsumer(
	consumerImpl nsq.Consumer,
	cc []nsq.ConsumerTopicChannel,
) error {
	for _, val := range cc {

		err := consumerImpl.GetConsumer(val.Topic, val.Channel, gnsq.NewConfig())

		if err != nil {
			log.Println("func NewConsumer", val.Topic, val.Channel, err)
			continue
		}

		consumerImpl.AddHandler(gnsq.HandlerFunc(val.HandlerFunc), val.Concurrency)

		err = consumerImpl.ConnectNSQLookupD()
		if err != nil {
			log.Println("func ConnectToNSQLookupd", val.Topic, val.Channel, err)
		}

		err = consumerImpl.ConnectNSQD()

		if err != nil {
			// it's mean once failed NSQD is die
			log.Println("failed to connect ", err)
			return err
		}
	}

	return nil
}

// Publish to NSQ consumer
func (g *gonsq) Publish(topic string, msg []byte) error {

	if err := g.prod.Publish(topic, msg); err != nil {
		log.Println("func Publish", err)
		return err
	}

	return nil
}
