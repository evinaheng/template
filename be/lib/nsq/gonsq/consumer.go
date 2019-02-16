package gonsq

import (
	"errors"

	nsqlib "github.com/nsqio/go-nsq"
	"github.com/template/be/lib/nsq"
)

type consumer struct {
	nsqc              *nsqlib.Consumer
	config            *nsqlib.Config
	NSQDAddress       string
	NSQLookupDAddress string
}

// NewConsumerMethodNSQ is get instance nsqlib consumer
func NewConsumerMethodNSQ(NSQDAddress string, NSQLookupDAddress string) (nsq.Consumer, error) {
	var err error

	return &consumer{
		config:            nsqlib.NewConfig(),
		NSQDAddress:       NSQDAddress,
		NSQLookupDAddress: NSQLookupDAddress,
	}, err
}

func (c *consumer) GetConsumer(topic, channel string, config *nsqlib.Config) (err error) {
	if config == nil {
		err = errors.New("config cannot be nil")
	} else {
		c.nsqc, err = nsqlib.NewConsumer(topic, channel, config)
	}

	return
}

func (c *consumer) ConnectNSQLookupD() error {
	return c.nsqc.ConnectToNSQLookupd(c.NSQLookupDAddress)
}
func (c *consumer) ConnectNSQD() error {
	return c.nsqc.ConnectToNSQD(c.NSQDAddress)
}
func (c *consumer) AddHandler(handler nsqlib.Handler, concurrency int) {
	c.nsqc.AddConcurrentHandlers(handler, concurrency)
}
