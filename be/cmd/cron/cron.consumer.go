package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	nsq "github.com/nsqio/go-nsq"
	c "github.com/robfig/cron"
	"github.com/template/be/cmd/internal"
	testCon "github.com/template/be/consumer/test"
	"github.com/template/be/lib/panics"
)

var allConsumers = []*nsq.Consumer{}

// Initialize all consumers
func initConsumer(panicWrapper panics.Panics, config internal.Config) {

	// Connect the consumer into lookupd endpoint from config
	nSQLookupds := strings.Split(config.NSQLookupds.Endpoint, ",")
	if len(nSQLookupds) < 1 {
		log.Fatalln("[FATAL] Wrong NSQLookupd Endpoint", config.NSQLookupds.Endpoint)
	}

	//---------------------
	// Consumer list
	//---------------------

	// - Test Example Consumer for Dummy
	fcdConfig := nsq.NewConfig()
	fcdConfig.MaxAttempts = 2
	fcdConfig.MaxInFlight = 3
	openConsumer(internal.ConsumerParam{
		Topic:     "test",
		Channel:   "test",
		LogPrefix: "nsq-test",
		Count:     1,
		Config:    fcdConfig,
		Handler:   panicWrapper.NSQ(testCon.TestNSQ),
	})

}

// Initialize consumer
func openConsumer(param internal.ConsumerParam) {

	// Use default config
	if param.Config == nil {
		param.Config = nsq.NewConfig()
	}

	// Use default name
	if param.Name == "" {
		param.Name = param.Topic + " " + param.Channel
	}

	consumers := make([]*nsq.Consumer, param.Count)
	for x := 0; x < param.Count; x++ {
		consumers[x] = startConsumer(param.Topic, param.Channel, param.LogPrefix+":", param.Config, param.Handler, param.NSQLookupds)
		if consumers[x] == nil {
			log.Println("[WARN] Can't connect to NSQ", param.Name)
		}
		allConsumers = append(allConsumers, consumers[x])
	}

	// Check for NSQ last hit using cron -- every 3 minutes
	var reconnectCron = c.New()
	reconnectCron.AddFunc("@every 3m", func() {

		for x := 0; x < param.Count; x++ {

			if consumers[x] == nil {

				// Consumer isn't initialized

				consumers[x] = startConsumer(param.Topic, param.Channel, param.LogPrefix, param.Config, param.Handler, param.NSQLookupds)

			} else if consumers[x].IsStarved() {

				// Reconnect to prevent IO error

				fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "Reconnect NSQ", param.Name)

				// Stop the current consumer and reconnect
				consumers[x].Stop()
				consumers[x] = startConsumer(param.Topic, param.Channel, param.LogPrefix, param.Config, param.Handler, param.NSQLookupds)
			}

		}

	})
	reconnectCron.Start()
}

// Start consumer connection to NSQ
func startConsumer(topic, channel, logPrefix string, config *nsq.Config, f nsq.HandlerFunc, nsqLookupds []string) *nsq.Consumer {

	// Create new consumer
	consumer, _ := nsq.NewConsumer(topic, channel, config)
	consumer.SetLogger(log.New(os.Stderr, logPrefix, log.Ltime), nsq.LogLevelError)
	consumer.AddHandler(f)

	// Open connection to NSQ
	if err := consumer.ConnectToNSQLookupds(nsqLookupds); err != nil {
		log.Println("[WARN] Can't connect to NSQ", err)
	}

	return consumer
}

// Stop all consumers
func stopConsumers() {
	for c := range allConsumers {
		if allConsumers[c] != nil {
			allConsumers[c].Stop()
		}
	}
}
