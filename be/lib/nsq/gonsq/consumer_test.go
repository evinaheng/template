package gonsq

import (
	"testing"

	"github.com/stretchr/testify/assert"

	nsq "github.com/nsqio/go-nsq"
)

func TestConsumer(t *testing.T) {
	stopChan := make(chan int)
	csm := consumer{
		nsqc: &nsq.Consumer{
			StopChan: stopChan,
		},
	}

	err := csm.ConnectNSQD()
	assert.NotNil(t, err)
	err = csm.ConnectNSQLookupD()
	assert.NotNil(t, err)

	_, err = NewConsumerMethodNSQ("some-nsqd", "some-nsqlookupd")

	assert.Nil(t, err)

	err = csm.GetConsumer("some-topic", "soime-channel", nil)
	assert.NotNil(t, err)
}
