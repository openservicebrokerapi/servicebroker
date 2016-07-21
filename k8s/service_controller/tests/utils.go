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
var serverCmd *exec.Cmd
var stdout bytes.Buffer
var stderr bytes.Buffer

func StartServer() error {
	// Make sure some other server isn't running
	if _, err := ServerGET("/v2/service_brokers"); err == nil {
		return fmt.Errorf("Server already running - stop it!")
	}

	// Should look into running the docker image instead
	serverCmd = exec.Command("../service_controller")
	serverCmd.Stdout = &stdout
	serverCmd.Stderr = &stderr
	err := serverCmd.Start()
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
	if serverCmd != nil {
		serverCmd.Process.Kill()
		serverCmd.Wait()
		serverCmd = nil
	}
	stdout.Reset()
	stderr.Reset()
}

func BrokerGET(url string) (string, error) {
	resp, err := http.Get(brokerURL + url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

var brokerURL = "http://localhost:9090/"
var brokerCmd *exec.Cmd
var brokerStdout bytes.Buffer
var brokerStderr bytes.Buffer

func StartBroker() error {
	// Make sure some other broker isn't running
	if _, err := BrokerGET("/v2/catalog"); err == nil {
		return fmt.Errorf("Broker already running - stop it!")
	}

	// Should look into running the docker image instead
	brokerCmd = exec.Command("../brokers/go/gobroker")
	brokerCmd.Stdout = &brokerStdout
	brokerCmd.Stderr = &brokerStderr
	err := brokerCmd.Start()
	if err != nil {
		StopBroker() // cleanup
		return err
	}

	// Wait for the broker to be available
	start := time.Now().Unix()
	for ; time.Now().Unix()-start < 10; time.Sleep(1 * time.Second) {
		if _, err = BrokerGET("/v2/catalog"); err == nil {
			return nil
		}
	}

	// Return last 'err' we got - it might help debug
	StopBroker() // cleanup
	return fmt.Errorf("Timed-out waiting for the broker: %s", err)
}

func StopBroker() {
	if brokerCmd != nil {
		brokerCmd.Process.Kill()
		brokerCmd.Wait()
		brokerCmd = nil
	}
	stdout.Reset()
	stderr.Reset()
}
