package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func MaskFields(t *testing.T, input interface{}, paths []string) interface{} {
	for _, path := range paths {
		input = MaskField(t, input, path)
	}
	return input
}

func MaskField(t *testing.T, input interface{}, path string) interface{} {
	var iface interface{}

	if str, ok := input.(string); ok {
		if err := json.Unmarshal([]byte(str), &iface); err != nil {
			t.Fatal(fmt.Sprintf("Can't unmarshal json: %s", str))
		}
	} else {
		iface = input
	}

	path = strings.TrimSpace(path)
	if path != "" {
		fields := strings.Split(path, ".")
		if ptr, ok := iface.(map[string]interface{}); ok {
			for i, field := range fields {
				if ptr[field] == nil {
					break
				}
				if i+1 != len(fields) {
					ptr = ptr[field].(map[string]interface{})
					continue
				}
				ptr[field] = "XXX"
			}
		} else if array, ok := iface.([]interface{}); ok {
			for _, item := range array {
				ptr := item.(map[string]interface{})
				for i, field := range fields {
					if ptr[field] == nil {
						break
					}
					if i+1 != len(fields) {
						ptr = ptr[field].(map[string]interface{})
						continue
					}
					ptr[field] = "XXX"
				}
			}
		}
	}
	return iface
}

func GetField(t *testing.T, inJson string, path string) string {
	var iface interface{}
	if err := json.Unmarshal([]byte(inJson), &iface); err != nil {
		t.Fatal(fmt.Sprintf("Can't unmarshal json: %s", inJson))
	}

	path = strings.TrimSpace(path)
	if path == "" {
		t.Fatal("GetField can't take an empty string")
	}

	fields := strings.Split(path, ".")
	ptr := iface.(map[string]interface{})
	for i, field := range fields {
		if i+1 != len(fields) {
			ptr = ptr[field].(map[string]interface{})
			continue
		}
		var s string
		var ok bool
		if s, ok = ptr[field].(string); !ok {
			t.Fatal("Can't find or convert field value into a string: " + field)
		}
		return s
	}
	t.Fatal("Should never get here")
	return ""
}

func ServerGET(url string) (string, *http.Response, error) {
	resp, err := http.Get(serverURL + url)
	if err != nil {
		return "", nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), resp, err
}

func ServerPOST(url string, cType string, data io.Reader) (string, *http.Response, error) {
	resp, err := http.Post(serverURL+url, cType, data)
	if err != nil {
		return "", nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), resp, err
}

var serverURL = "http://localhost:10000"
var serverCmd *exec.Cmd
var stdout bytes.Buffer
var stderr bytes.Buffer

func StartServer() error {
	// Make sure some other server isn't running
	if _, _, err := ServerGET("/v2/service_brokers"); err == nil {
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
		if _, _, err = ServerGET("/v2/service_brokers"); err == nil {
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

var brokerURL = "http://localhost:9090"
var brokerCmd *exec.Cmd
var brokerStdout bytes.Buffer
var brokerStderr bytes.Buffer

func StartBroker() error {
	// Make sure some other broker isn't running
	if _, err := BrokerGET("/v2/catalog"); err == nil {
		return fmt.Errorf("Broker already running - stop it!")
	}

	// Should look into running the docker image instead
	brokerCmd = exec.Command("../../../brokers/go/gobroker")
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
