package main

import (
	"testing"
	"os"
)

func TestMain(m *testing.M) {

	configFile := os.Getenv("CONFIG_FILE")

	if configFile == "" {
		configFile = "./tests/config.test.yml"
	}

	loadConfig(configFile)

	os.Exit(m.Run())
}
