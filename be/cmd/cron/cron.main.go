/*Package main

The entrance for Cron instance. Steps for initialization:
- Set logging flags to use long version (include filename and line)
- Check for testing mode. Process will end if it's TRUE
- Add GOPS listener
- Read config files
- Set global value for internal/global/global.go
- Load localization module
- Load resources (datadog, featureflag, iris)
- Load panics module
- Load usecases
- Start cronjob

*/

package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	gops "github.com/google/gops/agent"
	"github.com/template/be/cmd/internal"
	logging "github.com/template/be/lib/logging"
	"github.com/template/be/lib/server"
	"github.com/template/be/locale"
)

func main() {

	// Set logging flags
	logging.LogInit()
	log.SetFlags(log.LstdFlags | log.Llongfile)

	// Check config testing mode
	var isConfigTest bool
	flag.BoolVar(&isConfigTest, "test", false, "Enable config test mode")
	flag.Parse()
	if isConfigTest {
		internal.TestConfig()
	}

	// GOPS
	if err := gops.Listen(gops.Options{
		ShutdownCleanup: true,
	}); err != nil {
		log.Fatal("[FATAL] Can't initialize GOPS", err)
	}

	log.Println("Starting : API Cron")

	// Read config
	config := internal.InitConfig()

	// Initialize single httpClient
	/*httpClient := &http.Client{
		Timeout: time.Duration(config.Environment.GlobalTimeout) * time.Second,
	}*/

	// Get server IPAddress
	ipAddress := server.GetIPAddress()

	// Init localization
	locale.Init(config.Server.Env)

	// Init random
	rand.Seed(time.Now().UnixNano())

	// Initialize panic handling
	panicWrapper := internal.InitPanics(config.Server.Env, ipAddress)

	// Get all available usecases
	ucase := internal.GetUsecase(&config, ipAddress)

	// Set cron job
	initCron(panicWrapper, ucase)

	// Initialize consumer
	initConsumer(panicWrapper, config)

	log.Println("Running : API Cron")

	// Serve cron without HTTP
	serve()

}
