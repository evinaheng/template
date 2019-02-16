package toslack

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	nsq "github.com/nsqio/go-nsq"
	"github.com/template/be/lib/generate"
	"github.com/template/be/lib/panics"
	"github.com/template/be/lib/router"
	"github.com/template/be/lib/slack"
)

var newOnce sync.Once

var instance *panicsToSlack

// InitPanic module
// Can only init once
func InitPanic(config Config) panics.Panics {

	if config.SlackURL == "" {
		return nil
	}

	newOnce.Do(func() {

		instance = &panicsToSlack{
			env:         strings.ToUpper(config.Env),
			withMention: config.WithMention,
			slackURL:    config.SlackURL,
			ipAddress:   config.IPAddress,
		}

		go instance.captureBadDeployment()

	})

	return instance
}

// Handler for http middleware
func (p *panicsToSlack) HTTPHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestDump, _ := httputil.DumpRequest(req, true)
		defer func() {

			var err error
			r := recover()

			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}

				stack := "Request:\n" + string(requestDump) + "Stack:\n" + string(debug.Stack())

				// Publish to Slack
				p.publishError("HTTP", err.Error(), stack)
			}

		}()
		next.ServeHTTP(w, req)
	})

}

// RouterHandler for handling router function call
func (p *panicsToSlack) RouterHandler(ctx context.Context, request *http.Request, param *router.HandlerParam, f func(context.Context, *router.HandlerParam) router.HandlerResult) router.HandlerResult {

	// Prevent panic error
	defer func() {
		requestDump, _ := httputil.DumpRequest(request, true)

		var err error
		r := recover()

		if r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = errors.New("Unknown error")
			}

			stack := "Request:\n" + string(requestDump) + "Stack:\n" + string(debug.Stack())

			// Publish to Slack
			p.publishError("API", err.Error(), stack)

		}
	}()

	return f(ctx, param)

}

// Cron panics function wrapper
func (p *panicsToSlack) Cron(f func()) func() {

	return func() {
		// Prevent panic error
		defer func() {
			var err error
			r := recover()

			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}

				stack := "Stack:\n" + string(debug.Stack())

				// Publish to Slack
				p.publishError("Cron", err.Error(), stack)
			}
		}()

		f()
	}

}

// NSQ consumer panics function wrapper
func (p *panicsToSlack) NSQ(handler nsq.HandlerFunc) nsq.HandlerFunc {

	return func(message *nsq.Message) error {
		// Prevent panic error
		defer func() {
			var err error
			r := recover()

			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				var body string
				if message != nil {
					body = string(message.Body)
				}

				stack := "Body:\n" + body + "Stack:\n" + string(debug.Stack())

				// Publish to Slack
				p.publishError("NSQ", err.Error(), stack)
			}
		}()

		return handler(message)
	}

}

func (p *panicsToSlack) captureBadDeployment() {
	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGUSR1)
	for {
		select {
		case <-term:
			p.publishError("Deploy", "Failed to deploy", "")
		}
	}

}

func (p *panicsToSlack) publishError(handlerType, title, stack string) {

	// Generate unique ID
	id := generateID()

	// Print to error log
	log.Printf("PANICS! [%s] %s: %s \n", id, handlerType, title)

	// Generate attachment
	at := slack.Attachment{
		Color:     "#ff3f3f",
		TitleLink: fmt.Sprintf("https://www.scalyr.com/events?filter=%s&mode=log&logSource=%s", id, p.ipAddress),
		Title:     fmt.Sprintf("%s - [%s] %s: %s", id, p.env, handlerType, title),
		Text:      stack,
	}

	slack.SendAttachment(at, p.slackURL, p.withMention)
}

// Get random id for better readability
func generateID() string {
	id := generate.MD5(time.Now().String() + generate.RandomString(6, generate.StringAlphaNumeric))
	return id[:6]
}
