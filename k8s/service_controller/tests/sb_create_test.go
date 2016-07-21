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
	res, err := ServerGET("/v2/service_brokers")
	if err != nil {
		t.Fatal("Error getting list of brokers: %s", err)
	}
	if strings.TrimSpace(res) != "[]" {
		t.Fatal("Should be an empty list: %v", res)
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
	res, err := ServerPOST("/v2/service_brokers", "application/json", body)
	if err != nil {
		t.Fatal("Error registering a new broker: %s", err)
	}
	expected := `{
	  "metadata": {
	    "guid":"123-123",
		"created_at":"",
		"updated_at":"",
		"url":""
	  },
	  "entity":{
	    "name":"service-broker-name",
		"broker_url":"",
		"auth_username":""
      }
	}`

	maskFields := []string{
		"metadata.guid",
		"metadata.created_at",
		"metadata.updated_at",
	}

	d1 := MaskFields(t, expected, maskFields)
	d2 := MaskFields(t, res, maskFields)

	if !reflect.DeepEqual(d1, d2) {
		t.Fatal(fmt.Sprintf("Wrong results.\nExpected: %q\nGot: %q\n", d1, d2))
	}
}
