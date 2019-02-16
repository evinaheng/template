package toslack_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/h2non/gock"
	nsq "github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	. "github.com/template/be/lib/panics/toslack"
	"github.com/template/be/lib/router"
)

func TestInitPanicError(t *testing.T) {
	c := Config{}

	p := InitPanic(c)
	assert.Nil(t, p)
}

func TestInitPanic(t *testing.T) {
	c := Config{
		SlackURL: "https://www.slack.com",
	}

	p := InitPanic(c)
	assert.NotNil(t, p)
	time.Sleep(1 * time.Millisecond)
}

func TestHTTPHandler(t *testing.T) {

	defer gock.Off()
	gock.New("https://www.slack.com").Post("").MatchType("application/x-www-form-urlencoded").Reply(200).BodyString(`ok`)

	c := Config{
		SlackURL:  "https://www.slack.com",
		Env:       "unittest",
		IPAddress: "127.0.0.1",
	}

	p := InitPanic(c)
	assert.NotNil(t, p)

	dummyHandler1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aa := 0
		log.Println(1 / aa)
	})

	dummyHandler2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("foo")
	})

	dummyHandler3 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(123)
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()

	// Error object panic
	wrappedHandler := func() {
		p.HTTPHandler(dummyHandler1).ServeHTTP(w, req)
	}
	assert.NotPanics(t, wrappedHandler)

	// String panic
	wrappedHandler = func() {
		p.HTTPHandler(dummyHandler2).ServeHTTP(w, req)
	}
	assert.NotPanics(t, wrappedHandler)

	// Others panic
	wrappedHandler = func() {
		p.HTTPHandler(dummyHandler3).ServeHTTP(w, req)
	}
	assert.NotPanics(t, wrappedHandler)

}

func TestRouterHandler(t *testing.T) {
	defer gock.Off()
	gock.New("https://www.slack.com").Post("").MatchType("application/x-www-form-urlencoded").Reply(200).BodyString(`ok`)

	c := Config{
		SlackURL:  "https://www.slack.com",
		Env:       "unittest",
		IPAddress: "127.0.0.1",
	}

	p := InitPanic(c)
	assert.NotNil(t, p)

	dummyFunc1 := func(context.Context, *router.HandlerParam) router.HandlerResult {
		aa := 0
		log.Println(1 / aa)
		return router.HandlerResult{}
	}

	dummyFunc2 := func(context.Context, *router.HandlerParam) router.HandlerResult {
		panic("foo")
		return router.HandlerResult{}
	}

	dummyFunc3 := func(context.Context, *router.HandlerParam) router.HandlerResult {
		panic(123)
		return router.HandlerResult{}
	}

	req := httptest.NewRequest("GET", "/api/users", nil)

	// Error object panic
	wrappedHandler := func() {
		p.RouterHandler(context.TODO(), req, &router.HandlerParam{}, dummyFunc1)
	}
	assert.NotPanics(t, wrappedHandler)

	// Error string panic
	wrappedHandler = func() {
		p.RouterHandler(context.TODO(), req, &router.HandlerParam{}, dummyFunc2)
	}
	assert.NotPanics(t, wrappedHandler)

	// Others panic
	wrappedHandler = func() {
		p.RouterHandler(context.TODO(), req, &router.HandlerParam{}, dummyFunc3)
	}
	assert.NotPanics(t, wrappedHandler)

}

func TestCron(t *testing.T) {
	defer gock.Off()
	gock.New("https://www.slack.com").Post("").MatchType("application/x-www-form-urlencoded").Reply(200).BodyString(`ok`)

	c := Config{
		SlackURL:  "https://www.slack.com",
		Env:       "unittest",
		IPAddress: "127.0.0.1",
	}

	p := InitPanic(c)
	assert.NotNil(t, p)

	dummyFunc1 := func() {
		aa := 0
		log.Println(1 / aa)
	}

	dummyFunc2 := func() {
		panic("foo")
	}

	dummyFunc3 := func() {
		panic(123)
	}

	// Error object panic
	wrappedHandler := func() {
		p.Cron(dummyFunc1)()
	}
	assert.NotPanics(t, wrappedHandler)

	// Error string panic
	wrappedHandler = func() {
		p.Cron(dummyFunc2)()
	}
	assert.NotPanics(t, wrappedHandler)

	// Others panic
	wrappedHandler = func() {
		p.Cron(dummyFunc3)()
	}
	assert.NotPanics(t, wrappedHandler)

}

func TestNSQ(t *testing.T) {
	defer gock.Off()
	gock.New("https://www.slack.com").Post("").MatchType("application/x-www-form-urlencoded").Reply(200).BodyString(`ok`)

	c := Config{
		SlackURL:  "https://www.slack.com",
		Env:       "unittest",
		IPAddress: "127.0.0.1",
	}

	p := InitPanic(c)
	assert.NotNil(t, p)

	dummyFunc1 := nsq.HandlerFunc(func(msg *nsq.Message) error {
		aa := 0
		log.Println(1 / aa)
		return nil
	})

	dummyFunc2 := nsq.HandlerFunc(func(msg *nsq.Message) error {
		panic("foo")
		return nil
	})

	dummyFunc3 := nsq.HandlerFunc(func(msg *nsq.Message) error {
		panic(123)
		return nil
	})

	// Error object panic
	wrappedHandler := func() {
		p.NSQ(dummyFunc1).HandleMessage(&nsq.Message{})
	}
	assert.NotPanics(t, wrappedHandler)

	// Error string panic
	wrappedHandler = func() {
		p.NSQ(dummyFunc2).HandleMessage(nil)
	}
	assert.NotPanics(t, wrappedHandler)

	// Others panic
	wrappedHandler = func() {
		p.NSQ(dummyFunc3).HandleMessage(nil)
	}
	assert.NotPanics(t, wrappedHandler)

}
