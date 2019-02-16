package tolog

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"runtime/debug"

	nsq "github.com/nsqio/go-nsq"
	"github.com/template/be/lib/panics"
	"github.com/template/be/lib/router"
)

// New dummy panics
func New() panics.Panics {
	return &panicsToLog{}
}

func (p *panicsToLog) HTTPHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestDump, _ := httputil.DumpRequest(r, true)
		defer func() {

			r := recover()

			if r != nil {

				stack := "Request:\n" + string(requestDump) + "Stack:\n" + string(debug.Stack())

				log.Println(r)
				log.Println(stack)

			}

		}()
		next.ServeHTTP(w, r)
	})
}

// RouterHandler for handling router function call
func (p *panicsToLog) RouterHandler(ctx context.Context, request *http.Request, param *router.HandlerParam, f func(context.Context, *router.HandlerParam) router.HandlerResult) router.HandlerResult {

	// Prevent panic error
	defer func() {
		requestDump, _ := httputil.DumpRequest(request, true)

		r := recover()

		if r != nil {

			stack := "Request:\n" + string(requestDump) + "Stack:\n" + string(debug.Stack())

			log.Println(r)
			log.Println(stack)
		}
	}()

	return f(ctx, param)

}

// Cron panics function wrapper
func (p *panicsToLog) Cron(f func()) func() {

	return func() {
		// Prevent panic error
		defer func() {
			r := recover()

			if r != nil {

				stack := "Stack:\n" + string(debug.Stack())

				log.Println(r)
				log.Println(stack)
			}
		}()

		f()
	}

}

// NSQ consumer panics function wrapper
func (p *panicsToLog) NSQ(handler nsq.HandlerFunc) nsq.HandlerFunc {

	return func(message *nsq.Message) error {
		// Prevent panic error
		defer func() {
			r := recover()

			if r != nil {
				var body string
				if message != nil {
					body = string(message.Body)
				}

				stack := "Body:\n" + body + "Stack:\n" + string(debug.Stack())

				log.Println(r)
				log.Println(stack)
			}
		}()

		return handler(message)
	}

}
