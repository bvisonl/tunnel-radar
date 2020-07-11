package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	configFile := os.Getenv("CONFIG_FILE")

	if configFile == "" {
		configFile = "./tests/config.test.yml"
	}

	loadConfig(configFile)

	go StartCliServer(tunnelRadarConfig.CliServerHost, tunnelRadarConfig.CliServerPort)

	os.Exit(m.Run())
}
