package dummy

import (
	"errors"

	"github.com/template/be/lib/nsq"
)

// New GoNSQ module
func New(alwaysError bool) nsq.Producer {

	g := &dummynsq{
		alwaysError: alwaysError,
	}

	return g
}

// Publish to NSQ consumer
func (g *dummynsq) Publish(topic string, msg []byte) error {

	if g.alwaysError {
		return errors.New("Always error NSQ")
	}

	return nil
}
