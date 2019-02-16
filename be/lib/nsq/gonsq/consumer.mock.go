package gonsq

import (
	nsqlib "github.com/nsqio/go-nsq"
	"github.com/template/be/lib/nsq"
)

type mockConsumer struct {
}

// NewMockConsumer is mock for consumer
func NewMockConsumer() nsq.Consumer {
	return &mockConsumer{}
}

func (m *mockConsumer) GetConsumer(topic, channel string, config *nsqlib.Config) (err error) {
	return nil
}

func (m *mockConsumer) ConnectNSQLookupD() error {
	return nil
}
func (m *mockConsumer) ConnectNSQD() error {
	return nil
}

func (m *mockConsumer) AddHandler(handler nsqlib.Handler, concurrency int) {
}
