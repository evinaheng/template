package main

import (
	"log"
	"os"
	"os/signal"
)

func serve() {

	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan // wait for SIGINT

	stopConsumers()

	log.Println("Gracfully Stopped : API Cron")

}
