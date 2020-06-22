package main

import (
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

var tunnelRadarConfig TunnelRadarConfig

type TunnelRadarConfig struct {
	Tunnels map[string]*Tunnel `yaml:"tunnels"`
}

// LoadConfig - Load the configuration from ./config.yml
func loadConfig(configPath string) {

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&tunnelRadarConfig)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	// Configure clients
	for alias, tunnel := range tunnelRadarConfig.Tunnels {

		tunnel.Alias = alias

		// Set the Authentication method
		var authMethod ssh.AuthMethod
		if tunnel.Auth.Password != "" {
			authMethod = ssh.Password(tunnel.Auth.Password)
		} else if tunnel.Auth.Key != "" {
			buffer, err := ioutil.ReadFile(tunnel.Auth.Key)
			if err != nil {
				log.Fatalf("An error occurred loading the SSH key %s. Error: %s\r\n", tunnel.Auth.Key, err)
			}
			key, err := ssh.ParsePrivateKey(buffer)
			if err != nil {
				log.Fatalf("An error occurred loading the SSH key %s. Error: %s\r\n", tunnel.Auth.Key, err)
			}
			authMethod = ssh.PublicKeys(key)
		} else {
			log.Fatalf("You must specify an authentication method for %s\r\n", tunnel.Alias)
		}

		// Set the ClientConfig
		tunnel.ClientConfig = &ssh.ClientConfig{
			User:            tunnel.Auth.User,
			Auth:            []ssh.AuthMethod{authMethod},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

	}
}
