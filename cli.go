package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

const (
	TUNNEL_RADAR_EOF = "TUNNEL_RADAR_EOF\n"
)

func StartCli(host string, port int) {
	conn, err := net.Dial("tcp4", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalln(err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {

		// Send Command
		fmt.Printf("%s:%d> ", host, port)
		cmd, _ := reader.ReadString('\n')
		if cmd == "" {
			continue
		}
		conn.Write([]byte(cmd + "\n"))

		// Receive command
		reader := bufio.NewReader(conn)
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

}

func StartCliServer(host string, port int) {

	listen, err := net.Listen("tcp4", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalln(err)
	}

	defer listen.Close()

	fmt.Printf("Listening for CLI commands on %s\r\n", host+":"+strconv.Itoa(port))

	for {
		cliConn, err := listen.Accept()
		if err != nil {
			log.Printf("An error occurred accepting a connection from a CLI. Error: %s\r\n", err)
			continue
		}
		go getCommand(cliConn)
	}
}

func getCommand(conn net.Conn) {

	for {
		reader := bufio.NewReader(conn)
		data, _, err := reader.ReadLine()

		if err != nil && err != io.EOF {
			fmt.Println("read error:", err)
			return
		}

		if string(data) == "" {
			continue
		}

		cmd := strings.Split(string(data), " ")
		if len(cmd) <= 0 {
			conn.Write([]byte("Invalid command\n" + TUNNEL_RADAR_EOF))
			continue
		}

		switch strings.ToLower(cmd[0]) {
		case "list", "ls":
			list, err := tunnelRadarConfig.List()
			if err != nil {
				conn.Write([]byte(fmt.Sprintf("An error occurred getting the list. Error: %s\n%s", err, TUNNEL_RADAR_EOF)))
			} else {
				conn.Write([]byte(list + "\n" + TUNNEL_RADAR_EOF))
			}
			break
		case "enable":
			if len(cmd) <= 1 {
				conn.Write([]byte(fmt.Sprintf("Missing alias of tunnel (i.e. enable tunnel1)%s", TUNNEL_RADAR_EOF)))
			}
			alias := cmd[1]
			(*tunnelRadarConfig.Tunnels[alias]).Enable()
			conn.Write([]byte(alias + " has been enabled." + "\n" + TUNNEL_RADAR_EOF))

		case "disable":
			if len(cmd) <= 1 {
				conn.Write([]byte(fmt.Sprintf("Missing alias of tunnel (i.e. disable tunnel1)%s", TUNNEL_RADAR_EOF)))
			}
			alias := cmd[1]
			(*tunnelRadarConfig.Tunnels[alias]).Disable()
			conn.Write([]byte(alias + " has been disabled." + "\n" + TUNNEL_RADAR_EOF))

		case "exit":
			conn.Close()
			return

		default:
			conn.Write([]byte(fmt.Sprintf("Command %s not implemented yet.\r\n%s", data, TUNNEL_RADAR_EOF)))
		}

	}
}

func (config *TunnelRadarConfig) List() (list string, err error) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Alias", "Listen", "Remote", "Destination", "Status"})
	for alias, tunnel := range tunnelRadarConfig.Tunnels {
		t.AppendRow([]interface{}{alias, tunnel.Source, tunnel.Remote, tunnel.Destination, tunnel.Status})
	}
	t.SetStyle(table.StyleColoredBright)
	list = t.Render()

	return list, nil
}
