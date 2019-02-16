package tolog_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h2non/gock"
	nsq "github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	. "github.com/template/be/lib/panics/tolog"
	"github.com/template/be/lib/router"
)

func TestHandler(t *testing.T) {

	p := New()
	assert.NotNil(t, p)

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aa := 0
		log.Println(1 / aa)
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()

	wrappedHandler := func() {
		p.HTTPHandler(dummyHandler).ServeHTTP(w, req)
	}

	assert.NotPanics(t, wrappedHandler)

}

func TestRouterHandler(t *testing.T) {
	defer gock.Off()
	gock.New("https://www.slack.com").Post("").MatchType("application/x-www-form-urlencoded").Reply(200).BodyString(`ok`)

	p := New()
	assert.NotNil(t, p)

	dummyFunc := func(context.Context, *router.HandlerParam) router.HandlerResult {
		aa := 0
		log.Println(1 / aa)
		return router.HandlerResult{}
	}

	req := httptest.NewRequest("GET", "/api/users", nil)

	wrappedHandler := func() {
		p.RouterHandler(context.TODO(), req, &router.HandlerParam{}, dummyFunc)
	}

	assert.NotPanics(t, wrappedHandler)

}

func TestCron(t *testing.T) {
	defer gock.Off()
	gock.New("https://www.slack.com").Post("").MatchType("application/x-www-form-urlencoded").Reply(200).BodyString(`ok`)

	p := New()
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

	p := New()
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
