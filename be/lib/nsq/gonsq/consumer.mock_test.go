package gonsq

import (
	"testing"

	nsq "github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
)

func TestNewConsumerMock(t *testing.T) {
	mock := NewMockConsumer()

	err := mock.ConnectNSQLookupD()
	assert.Nil(t, err)

	err = mock.ConnectNSQD()
	assert.Nil(t, err)

	mock.GetConsumer("some-topic", "some-channel", &nsq.Config{
		Hostname: "some-hostname",
	})

}
