package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ghodss/yaml"
)

func callService(path, method, action string, reader io.ReadCloser) error {
	u := fmt.Sprintf("%s%s", controller, path)

	resp, err := callHttp(u, method, action, reader)
	if err != nil {
		return err
	}
	var j interface{}
	if err := json.Unmarshal([]byte(resp), &j); err != nil {
		return fmt.Errorf("Failed to parse JSON response from service: %s : %v", resp, err)
	}

	y, err := yaml.Marshal(j)
	if err != nil {
		return fmt.Errorf("Failed to serialize JSON response from service: %s : %v", resp, err)
	}

	fmt.Println(string(y))
	return nil
}

func callHttp(path, method, action string, reader io.ReadCloser) (string, error) {
	request, err := http.NewRequest(method, path, reader)
	request.Header.Add("Content-Type", "application/json")

	client := http.Client{
		Timeout: time.Duration(time.Duration(timeout) * time.Second),
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		message := fmt.Sprintf("status code: %d status: %s : %s", response.StatusCode, response.Status, body)
		return "", fmt.Errorf("cannot %s: %s\n", action, message)
	}

	return string(body), nil
}
