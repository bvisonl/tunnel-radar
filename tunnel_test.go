package main

import (
	"net/http"
	"testing"
)

func TestTunnel(t *testing.T) {

	wg.Wait() // wait for the CLI to finish testing

	// Make request
	resp, err := http.Get("http://127.0.0.1:3100")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Invalid response code when testing tunnel. Status Code: %d", resp.StatusCode)
	}
}
