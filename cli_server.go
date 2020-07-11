package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type CliServer struct {
	Host string
	Port int
	Conn net.Conn
}

func StartCliServer(host string, port int) {
	server := &CliServer{
		Host: host,
		Port: port,
	}
	server.listen()
}

func (server *CliServer) listen() {

	listen, err := net.Listen("tcp4", server.Host+":"+strconv.Itoa(server.Port))
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer listen.Close()

	fmt.Printf("Listening for CLI commands on %s\r\n", server.Host+":"+strconv.Itoa(server.Port))

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("An error occurred accepting a connection from a CLI. Error: %s\r\n", err)
			continue
		}
		go server.getCommands(conn)
	}
}

func (server *CliServer) getCommands(conn net.Conn) {

	for {
		reader := bufio.NewReader(conn)
		data, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			fmt.Println("read error:", err)
			return
		}
		result, err := server.handleCommand(string(data))
		conn.Write([]byte(fmt.Sprintf("%s\r\n%s", result, TUNNEL_RADAR_EOF)))
	}
}

func (server *CliServer) handleCommand(command string) (string, error) {

	cmd := strings.Split(command, " ")

	if len(cmd) <= 0 {
		return "", errors.New("Invalid command")
	}

	switch strings.ToLower(cmd[0]) {
	case "list", "ls":
		list, err := tunnelRadarConfig.list()
		if err != nil {
			return "", err
		} else {
			return list, nil
		}
		break
	case "enable":
		if len(cmd) <= 1 {
			return "", errors.New("Missing alias of tunnel (i.e. enable tunnel1)")
		}
		alias := cmd[1]
		(*tunnelRadarConfig.Tunnels[alias]).Enable()
		return fmt.Sprintf("%s has been enabled.", alias), nil

	case "disable":
		if len(cmd) <= 1 {
			return "", errors.New("Missing alias of tunnel (i.e. disable tunnel1)")
		}
		alias := cmd[1]
		(*tunnelRadarConfig.Tunnels[alias]).Disable()
		return fmt.Sprintf("%s has been disabled.", alias), nil

	case "exit":
		server.Conn.Close()
		return "Bye", nil

	default:
		return "", errors.New(fmt.Sprintf("Command %s not implemented yet.", cmd[0]))
	}

	return "", errors.New("Unknown error occurred")
}
