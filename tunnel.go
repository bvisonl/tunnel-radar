package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jedib0t/go-pretty/text"
	"golang.org/x/crypto/ssh"
)

type Tunnel struct {
	Alias  string `yaml:"alias"`
	Source string `yaml:"source"`
	Auth   struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Key      string `yaml:"key"`
	} `yaml:"auth"`
	Remote       string `yaml:"remote"`
	Destination  string `yaml:"destination"`
	Disabled     bool   `yaml:"disabled"`
	Status       string
	ClientConfig *ssh.ClientConfig
}

func (tunnel *Tunnel) Spawn() error {

	if tunnel.Disabled == true {
		tunnel.Status = text.FgRed.Sprint("OFFLINE")
		fmt.Println(tunnel)
		return nil
	}

	// Start accepting connections
	listener, err := net.Listen("tcp", tunnel.Source)
	if err != nil {
		log.Printf("An error occurred while connecting to %s. Error: %s\r\n", tunnel.Alias, err)
		return err
	}

	defer listener.Close()

	log.Printf("Tunnel %s active and listening on %s => %s => %s\r\n", tunnel.Alias, tunnel.Source, tunnel.Remote, tunnel.Destination)
	for {
		if tunnel.Disabled == true {
			listener.Close()
			tunnel.Status = text.FgRed.Sprint("OFFLINE")
			return nil
		}

		tunnel.Status = text.FgGreen.Sprint("ONLINE")

		aConn, err := listener.Accept()
		if err != nil {
			log.Printf("An error occurred while accepting a connection on %s. Error: %s\r\n", tunnel.Alias, err)
			return err
		}

		go tunnel.Flow(aConn)
	}

	return nil
}

func (tunnel *Tunnel) Flow(aConn net.Conn) error {

	bConn, err := ssh.Dial("tcp", tunnel.Remote, tunnel.ClientConfig)
	if err != nil {
		tunnel.Disabled = true
		log.Printf("An error occurred establishing the SSH connection with %s. Error: %s\r\n", tunnel.Alias, err)
	}

	cConn, err := bConn.Dial("tcp", tunnel.Destination)
	if err != nil {
		tunnel.Disabled = true
		log.Printf("An error occurred establishing the SSH connection with the destination of %s. Error: %s\r\n", tunnel.Alias, err)
	}

	go tunnel.ProxyData(aConn, cConn)
	go tunnel.ProxyData(cConn, aConn)

	return nil
}

func (tunnel *Tunnel) ProxyData(legA, legB net.Conn) {
	_, err := io.Copy(legA, legB)
	if err != nil {
		tunnel.Disabled = true
		log.Printf("An error occurred while forwarding connections between endpoints on %s. Error: %s\r\n", tunnel.Alias, err)
	}
}
