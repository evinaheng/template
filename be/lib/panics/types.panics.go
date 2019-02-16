package panics

import (
	"context"
	"net/http"

	nsq "github.com/nsqio/go-nsq"
	"github.com/template/be/lib/router"
)

// A Panics interface provides logging for panic situation
type Panics interface {
	HTTPHandler(http.Handler) http.Handler
	RouterHandler(context.Context, *http.Request, *router.HandlerParam, func(context.Context, *router.HandlerParam) router.HandlerResult) router.HandlerResult
	Cron(func()) func()
	NSQ(nsq.HandlerFunc) nsq.HandlerFunc
}
