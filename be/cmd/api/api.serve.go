package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/nytimes/gziphandler"
)

func serve(hport string, handler http.Handler) error {

	// Use 30 seconds timeout for server
	// We don't use global timeout because it will create race with API timeout
	var timeout = 30 * time.Second

	// error variable
	var err error
	// listener variable
	var l net.Listener

	fd := os.Getenv("EINHORN_FDS")
	if fd != "" {
		sock, err := strconv.Atoi(fd)
		if err == nil {
			hport = "socketmaster:" + fd
			file := os.NewFile(uintptr(sock), "listener")
			fl, err := net.FileListener(file)
			if err == nil {
				l = fl
			}
		}
	}

	if l == nil {
		var err error
		l, err = net.Listen("tcp4", hport)
		if err != nil {
			log.Fatal(err)
		}
	}

	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	srv := &http.Server{
		Addr:         hport,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
		IdleTimeout:  timeout,
		Handler:      gziphandler.GzipHandler(handler), // Set GZIP compression
	}

	// Set context timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go func() {
		// service connections
		if err = srv.Serve(l); err != nil {
			log.Println("func serve", err)
		}
	}()

	<-stopChan // wait for SIGINT
	log.Println("Stopping : API")

	// Shutdown the context
	srv.Shutdown(ctx)

	log.Println("Gracfully Stopped :  API")
	return err

}
