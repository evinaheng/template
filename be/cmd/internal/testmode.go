package internal

import (
	"log"
	"os"
)

// TestConfig runs diagnosis for connection
func TestConfig() {
	log.Println("Start Configuration Testing")

	// Read config
	log.Println("Reading config...")
	InitConfig()

	log.Println("Test successful!")
	os.Exit(0)
}
