package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

func ServerGET(url string) (string, error) {
	resp, err := http.Get(serverURL + url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

var serverURL = "http://localhost:10000/"
var cmd *exec.Cmd
var stdout bytes.Buffer
var stderr bytes.Buffer

func StartServer() error {
	// Make sure some other server isn't running
	if _, err := ServerGET("/v2/service_brokers"); err == nil {
		return fmt.Errorf("Server already running - stop it!")
	}

	// Should look into running the docker image instead
	cmd = exec.Command("../service_controller")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		StopServer() // cleanup
		return err
	}

	// Wait for the server to be available
	start := time.Now().Unix()
	for ; time.Now().Unix()-start < 10; time.Sleep(1 * time.Second) {
		if _, err = ServerGET("/v2/service_brokers"); err == nil {
			return nil
		}
	}

	// Return last 'err' we got - it might help debug
	StopServer() // cleanup
	return fmt.Errorf("Timed-out waiting for the server: %s", err)
}

func StopServer() {
	if cmd != nil {
		cmd.Process.Kill()
		cmd.Wait()
		cmd = nil
	}
	stdout.Reset()
	stderr.Reset()
}
