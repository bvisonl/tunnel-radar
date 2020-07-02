package main

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// TODO: Valdiate complete structure of the config file
	// Configuration loaded from TestMain
	if len(tunnelRadarConfig.Tunnels) <= 0 {
		t.Errorf("Unable to load configurations")
	}

}
