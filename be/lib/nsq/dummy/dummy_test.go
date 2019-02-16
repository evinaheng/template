package dummy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/template/be/lib/nsq/dummy"
)

func TestPublishSuccess(t *testing.T) {
	module := New(false)
	assert.Nil(t, module.Publish("foo", nil))
}
func TestPublishError(t *testing.T) {
	module := New(true)
	assert.EqualError(t, module.Publish("foo", nil), "Always error NSQ")
}
