package tests

import (
	"testing"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func TestMain(m *testing.M) {
	err := StartServer()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
	defer StopServer()
	os.Exit(m.Run())
}

func Test_SB_ping(t *testing.T) {
	_, err := ServerGET("/v2/service_brokers")
	if err != nil {
		t.Fatal(err)
	}
}

// Utils stuff
var serverURL = "http://localhost:10000/"

func ServerGET(url string) (string, error) {
	resp, err := http.Get(serverURL + url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

var cmd *exec.Cmd
var stdout bytes.Buffer
var stderr bytes.Buffer

func StartServer() error {
	cmd = exec.Command("../service_controller")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Wait for the server to be available
	for {
		if _, err = ServerGET("/v2/service_brokers"); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func StopServer() {
	cmd.Process.Kill()
	cmd.Wait()
}
