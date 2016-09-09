package tests

import (
	"testing"

	"fmt"
	"os"
	"reflect"
	"strings"
)

func TestMain(m *testing.M) {
	err := StartServer()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
	err = StartBroker()
	if err != nil {
		fmt.Printf("Error starting broker: %s\n", err)
		os.Exit(1)
	}
	rc := m.Run()
	StopBroker()
	StopServer()
	os.Exit(rc)
}

func Test_SB_Ping(t *testing.T) {
	body, resp, err := ServerGET("/v2/service_brokers")
	if err != nil {
		t.Fatalf("Error getting list of brokers: %s", err)
	}
	if resp.Header["Content-Type"][0] != "application/json" {
		t.Fatalf("Wrong Content-Type. Got %q expected 'application/json", resp.Header["Content-Type"])
	}
	if strings.TrimSpace(body) != "[]" {
		t.Fatalf("Should be an empty list: %v", body)
	}
}

func Test_SB_Create(t *testing.T) {
	data := `{
	  "name": "service-broker-name",
	  "broker_url": "http://localhost:9090",
	  "auth_username": "admin",
	  "auth_password": "secretpassw0rd"
	}`
	body := strings.NewReader(data)
	data, resp, err := ServerPOST("/v2/service_brokers", "application/json", body)
	if err != nil {
		t.Fatalf("Error registering a new broker: %s", err)
	}
	if resp.Header["Content-Type"][0] != "application/json" {
		t.Fatalf("Wrong Content-Type. Got %q expected 'application/json", resp.Header["Content-Type"])
	}
	expected := `{
	  "metadata": {
	    "guid":"",
		"created_at":"",
		"url":""
	  },
	  "entity":{
	    "name":"service-broker-name",
		"broker_url":"http://localhost:9090",
		"auth_username":"admin"
      }
	}`

	maskFields := []string{
		"metadata.guid",
		"metadata.created_at",
		"metadata.url",
	}

	d1 := MaskFields(t, expected, maskFields)
	d2 := MaskFields(t, data, maskFields)

	if !reflect.DeepEqual(d1, d2) {
		t.Fatal(fmt.Sprintf("Wrong results.\nExpected: %q\nGot: %q\n", d1, d2))
	}

	// Now make sure its there
	/*
		rbody, resp, err := ServerGET("/v2/service_brokers")
		if err != nil {
			t.Fatalf("Error getting list of brokers: %s", err)
		}
		if resp.Header["Content-Type"][0] != "application/json" {
			t.Fatalf("Wrong Content-Type. Got %q expected 'application/json", resp.Header["Content-Type"])
		}

		d1 = MaskFields(t, "["+expected+"]", maskFields)
		d2 = MaskFields(t, rbody, maskFields)

		if !reflect.DeepEqual(d1, d2) {
			t.Fatal(fmt.Sprintf("Wrong results.\nExpected: %q\nGot: %q\n", d1, d2))
		}
	*/
}
