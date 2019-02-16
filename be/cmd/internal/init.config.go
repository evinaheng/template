package internal

import (
	"io/ioutil"
	"log"
	"strings"

	gcfg "gopkg.in/gcfg.v1"
)

// InitConfig Read and process config file
func InitConfig() Config {

	appconfig := Config{}
	if ok := readConfig(&appconfig, "/etc/template") || readConfig(&appconfig, "./development"); !ok {
		log.Fatalln("Failed to read configuration file")
	}
	return appconfig
}

// readConfig is file handler for reading configuration files into variable
// Return: - boolean
func readConfig(ac *Config, path string) bool {
	parts := []string{"main", "db", "vendor"}
	var configString []string

	for _, v := range parts {
		fname := path + "/" + v + ".ini"
		config, err := ioutil.ReadFile(fname)
		if err != nil {
			return false
		}

		configString = append(configString, string(config))
	}

	if err := gcfg.ReadStringInto(ac, strings.Join(configString, "\n\n")); err != nil {
		log.Println("func readConfig", err)
		return false
	}

	return true
}
