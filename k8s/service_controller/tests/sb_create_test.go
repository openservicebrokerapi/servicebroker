package tests

import (
	"testing"

	"fmt"
	"os"
	"strings"
)

func TestMain(m *testing.M) {
	err := StartServer()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
	rc := m.Run()
	StopServer()
	os.Exit(rc)
}

func Test_SB_Ping(t *testing.T) {
	res, err := ServerGET("/v2/service_brokers")
	if err != nil {
		t.Fatal("Error getting list of brokers: %s", err)
	}
	if strings.TrimSpace(res) != "[]" {
		t.Fatal("Should be an empty list: %v", res)
	}
}
