package main

import (
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
	Listener     net.Listener
}

func (tunnel *Tunnel) Spawn() error {

	if tunnel.Listener != nil {
		tunnel.Listener.Close()
	}

	sshConn, err := ssh.Dial("tcp", tunnel.Remote, tunnel.ClientConfig)
	if err != nil {
		log.Printf("An error occurred establishing the SSH connection with %s. Error: %s\r\n", tunnel.Alias, err)
		tunnel.Disable()
		return err
	}

	// Start accepting connections
	listener, err := net.Listen("tcp", tunnel.Source)
	tunnel.Listener = listener
	if err != nil {
		log.Printf("An error occurred while connecting to %s. Error: %s\r\n", tunnel.Alias, err)
		return err
	}

	defer listener.Close()

	tunnel.Status = text.FgGreen.Sprint("ONLINE")

	log.Printf("Tunnel %s active and listening on %s => %s => %s\r\n", tunnel.Alias, tunnel.Source, tunnel.Remote, tunnel.Destination)
	for tunnel.Disabled == false {
		conn, err := listener.Accept()

		if err != nil && err != io.EOF {
			log.Printf("An error occurred while accepting a connection on %s. Error: %s\r\n", tunnel.Alias, err)
			return err
		}

		go tunnel.Flow(conn, sshConn)
	}
	return nil
}

func (tunnel *Tunnel) Flow(conn net.Conn, sshConn *ssh.Client) error {
	destConn, err := sshConn.Dial("tcp", tunnel.Destination)
	if err != nil {
		log.Printf("An error occurred establishing the SSH connection with the destination of %s. Error: %s\r\n", tunnel.Alias, err)
		tunnel.Disable()
		return err
	}

	go tunnel.ProxyData(conn, destConn)
	go tunnel.ProxyData(destConn, conn)

	return nil
}

func (tunnel *Tunnel) ProxyData(legA, legB net.Conn) {
	_, err := io.Copy(legA, legB)
	if err != nil {
		log.Printf("An error occurred while forwarding connections between endpoints on %s. Error: %s\r\n", tunnel.Alias, err)
	}
}

func (tunnel *Tunnel) Enable() {
	tunnel.Status = text.FgYellow.Sprint("CONNECTING")
	tunnel.Disabled = false
	go tunnel.Spawn()
}
func (tunnel *Tunnel) Disable() {
	tunnel.Status = text.FgRed.Sprint("OFFLINE")
	tunnel.Disabled = true
}
