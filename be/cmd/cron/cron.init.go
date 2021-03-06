package main

import (
	c "github.com/robfig/cron"
	"github.com/template/be/cmd/internal"
	"github.com/template/be/lib/panics"
)

var (
	cronObj *c.Cron
)

func initCron(panicWrapper panics.Panics, ucase *internal.Usecase) {

	// Initialize cron object
	cronObj = c.New()

	//---------------------
	// START - Cron list
	//---------------------

	// -- Open files logging datadog
	// -- Every 10 seconds
	cronObj.AddFunc("@every 10s", panicWrapper.Cron(func() {
		ucase.System.LogOpenFile()
	}))

	// END - Cron list
	//---------------------

	cronObj.Start()
}
