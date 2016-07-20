package server

import (
	"strings"
	"testing"

	"github.com/cncf/servicebroker/k8s/service_controller/model"
)

func TestNoBrokers(t *testing.T) {
	s := CreateInMemServiceStorage()
	l, err := s.ListBrokers()
	if err != nil {
		t.Fatalf("ListBrokers failed with: %#v", err)
	}
	if len(l) != 0 {
		t.Fatalf("Expected 0 brokers, got %d", len(l))
	}
	b, err := s.GetBroker("NOT THERE")
	if err == nil {
		t.Fatalf("GetBroker did not fail")
	}
	if b != nil {
		t.Fatalf("Got back a broker: %#v", b)
	}
}

func TestAddBroker(t *testing.T) {
	s := CreateInMemServiceStorage()
	b := &model.ServiceBroker{Name: "Test"}
	cat := model.Catalog{
		Services: []model.Service{},
	}
	err := s.AddBroker(b, &cat)
	if err != nil {
		t.Fatalf("AddBroker failed with: %#v", err)
	}
	l, err := s.ListBrokers()
	if len(l) != 1 {
		t.Fatalf("Expected 1 broker, got %d", len(l))
	}
	b2, err := s.GetBroker("Test")
	if err != nil {
		t.Fatalf("GetBroker failed: %#v", err)
	}
	if b2 == nil {
		t.Fatalf("Did not get back a broker")
	}
	if strings.Compare(b2.Name, b.Name) != 0 {
		t.Fatalf("Names don't match, expected: '%s', got '%s'", b.Name, b2.Name)
	}
}

func TestAddDuplicateBroker(t *testing.T) {
	s := CreateInMemServiceStorage()
	b := &model.ServiceBroker{Name: "Test"}
	cat := model.Catalog{
		Services: []model.Service{},
	}
	err := s.AddBroker(b, &cat)
	if err != nil {
		t.Fatalf("AddBroker failed with: %#v", err)
	}
	l, err := s.ListBrokers()
	if len(l) != 1 {
		t.Fatalf("Expected 1 broker, got %d", len(l))
	}
	b2, err := s.GetBroker("Test")
	if err != nil {
		t.Fatalf("GetBroker failed: %#v", err)
	}
	if b2 == nil {
		t.Fatalf("Did not get back a broker")
	}
	if strings.Compare(b2.Name, b.Name) != 0 {
		t.Fatalf("Names don't match, expected: '%s', got '%s'", b.Name, b2.Name)
	}
	err = s.AddBroker(b, &cat)
	if err == nil {
		t.Fatalf("AddBroker did not fail with duplicate")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("Unexpected error, wanted 'already exists' but got %#v", err)
	}
}
