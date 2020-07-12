package main

import (
	"testing"
	"time"
)

func TestCliClient(t *testing.T) {

	defer wg.Done()

	client, err := NewCliClient("127.0.0.1", 3101)
	if err != nil {
		t.Errorf(err.Error())
	}
	client.sendCommand("ls")
	time.Sleep(1 * time.Second)
	client.sendCommand("disable testServer")
	time.Sleep(1 * time.Second)
	client.sendCommand("enable testServer")
	time.Sleep(1 * time.Second)

}
