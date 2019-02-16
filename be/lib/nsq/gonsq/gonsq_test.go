package gonsq_test

import (
	"testing"

	gnsq "github.com/nsqio/go-nsq"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/template/be/lib/nsq"
	. "github.com/template/be/lib/nsq/gonsq"
)

func TestInitFailed(t *testing.T) {

	module, err := NewProducer(Config{
		Endpoint: "error",
	})

	assert.Nil(t, module)
	assert.EqualError(t, err, "Test error")

}

func TestInitSuccess(t *testing.T) {
	defer gock.Off()
	gock.New("tcp://10.164.4.31").Post("").Reply(200).BodyString("")

	module, _ := NewProducer(Config{
		Endpoint: "dummy:53",
	})

	assert.Error(t, module.Publish("foo", []byte("bar")))
}

func TestInitSuccessConsumer(t *testing.T) {

	mock := NewMockConsumer()

	listConsumer := make([]nsq.ConsumerTopicChannel, 0)

	listConsumer = append(listConsumer, nsq.ConsumerTopicChannel{
		Topic:   "some-topic",
		Channel: "some-channel",
		HandlerFunc: func(msg *gnsq.Message) error {
			return nil
		},
	})

	err := NewConsumer(mock, listConsumer)

	assert.Nil(t, err)

}
