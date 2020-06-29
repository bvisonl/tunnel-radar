package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	configFile := os.Getenv("CONFIG_FILE")

	loadConfig(configFile)

	if len(tunnelRadarConfig.Tunnels) <= 0 {
		t.Errorf("Unable to load configurations")
	}

}
