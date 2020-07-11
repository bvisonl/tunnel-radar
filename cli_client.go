package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type CliClient struct {
	Host   string
	Port   int
	Conn   net.Conn
	Reader *bufio.Reader
}

func NewCliClient(host string, port int) (*CliClient, error) {

	conn, err := net.Dial("tcp4", host+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(os.Stdin)

	cli := &CliClient{
		Host:   host,
		Port:   port,
		Conn:   conn,
		Reader: reader,
	}

	return cli, nil
}

func (client *CliClient) listen() {

	for {
		cmd, _ := client.Reader.ReadString('\n')
		if cmd == "" {
			continue
		}
		fmt.Printf("%s:%d> ", client.Host, client.Port)
		client.sendCommand(cmd)
	}
}

func (client *CliClient) sendCommand(cmd string) {
	// Send Command
	client.Conn.Write([]byte(cmd + "\n"))

	// Receive command
	reader := bufio.NewReader(client.Conn)
	var buffer bytes.Buffer

	for {
		data, err := reader.ReadBytes('\n')

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Print(err)
			break
		}

		if string(data) == TUNNEL_RADAR_EOF {
			break
		}

		buffer.Write(data)
	}

	fmt.Println(buffer.String())
}
