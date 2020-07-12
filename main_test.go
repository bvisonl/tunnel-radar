package main

import (
	"os"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func TestMain(m *testing.M) {

	configFile := os.Getenv("CONFIG_FILE")

	if configFile == "" {
		configFile = "./tests/config.test.yml"
	}

	loadConfig(configFile)

	tunnel := tunnelRadarConfig.Tunnels["testServer"]
	go tunnel.Spawn()

	go StartCliServer(tunnelRadarConfig.CliServerHost, tunnelRadarConfig.CliServerPort)

	// Wait for the tunnel and the cli server to start
	time.Sleep(1 * time.Second)

	wg.Add(1) // needed for now to wait for the CLI to finish testing

	os.Exit(m.Run())
}
