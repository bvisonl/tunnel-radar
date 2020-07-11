package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

const (
	TUNNEL_RADAR_EOF = "TUNNEL_RADAR_EOF\n"
)

func main() {

	time.Sleep(time.Second * 2)

	// Configurations by flag
	withCli := flag.Bool("i", false, "Enter in CLI mode to execute commands on TunnelRadar")
	cliHost := flag.String("ih", "127.0.0.1", "CLI host to connect to (default 127.0.0.1)")
	cliPort := flag.Int("ip", 7779, "CLI Port to connect to (default 7779)")

	configPath := flag.String("c", "/etc/tunnel-radar/config.yml", "Configuration file path")
	debug := flag.Bool("d", false, "Enable debugging")
	flag.Parse()

	if *withCli == true {
		cliClient, err := NewCliClient(*cliHost, *cliPort)
		if err != nil {
			log.Fatalf("An error occurred connecting to the CLI server. Error: %s\r\n", err)
		}

		// Bind the terminal and listen for commands
		defer cliClient.Conn.Close()
		cliClient.listen()
	} else {

		// Load configuration
		loadConfig(*configPath)

		if *debug == true {
			fmt.Println("Debug mode ON")
			fmt.Printf("Configuration loaded %s\r\n", *configPath)
		}

		for _, tunnel := range tunnelRadarConfig.Tunnels {
			go tunnel.Spawn()
		}

		// Start listening for commands
		StartCliServer(tunnelRadarConfig.CliServerHost, tunnelRadarConfig.CliServerPort)

	}
}
