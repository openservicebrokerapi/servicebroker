package mem

import (
	"strings"
	"testing"

	model "github.com/cncf/servicebroker/model/service_controller"
)

const (
	BROKER_UUID   = "126b8154-a24a-4e79-9185-3df2eb4d18a8"
	BROKER_UUID_2 = "2b0c42ed-c43a-4724-b883-e5ba878a8bfd"
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
	b := &model.ServiceBroker{GUID: "Test"}
	cat := model.Catalog{
		Services: []*model.Service{},
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
	b := &model.ServiceBroker{GUID: "Test"}
	cat := model.Catalog{
		Services: []*model.Service{},
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

func TestDeleteBroker(t *testing.T) {
	s := CreateInMemServiceStorage()
	b := &model.ServiceBroker{GUID: BROKER_UUID}
	cat := model.Catalog{
		Services: []*model.Service{},
	}
	err := s.AddBroker(b, &cat)
	if err != nil {
		t.Fatalf("AddBroker failed with: %#v", err)
	}
	l, err := s.ListBrokers()
	if len(l) != 1 {
		t.Fatalf("Expected 1 broker, got %d", len(l))
	}
	b2, err := s.GetBroker(BROKER_UUID)
	if err != nil {
		t.Fatalf("GetBroker failed: %#v", err)
	}
	if b2 == nil {
		t.Fatalf("Did not get back a broker")
	}
	if strings.Compare(b2.Name, b.Name) != 0 {
		t.Fatalf("Names don't match, expected: '%s', got '%s'", b.Name, b2.Name)
	}
	err = s.DeleteBroker(BROKER_UUID)
	if err != nil {
		t.Fatalf("Failed to delete broker: %s : %#v", BROKER_UUID, err)
	}
	l, err = s.ListBrokers()
	if len(l) != 0 {
		t.Fatalf("Expected 0 broker, got %d", len(l))
	}
	b2, err = s.GetBroker(BROKER_UUID)
	if err == nil {
		t.Fatalf("GetBroker returned a broker when there should be none")
	}
}

func TestDeleteBrokerMultiple(t *testing.T) {
	s := CreateInMemServiceStorage()
	b := &model.ServiceBroker{GUID: BROKER_UUID}
	b2 := &model.ServiceBroker{GUID: BROKER_UUID_2}
	cat := model.Catalog{
		Services: []*model.Service{{Name: "first"}},
	}
	cat2 := model.Catalog{
		Services: []*model.Service{{Name: "second"}},
	}
	err := s.AddBroker(b, &cat)
	if err != nil {
		t.Fatalf("AddBroker failed with: %#v", err)
	}
	err = s.AddBroker(b2, &cat2)
	if err != nil {
		t.Fatalf("AddBroker failed with: %#v", err)
	}
	l, err := s.ListBrokers()
	if len(l) != 2 {
		t.Fatalf("Expected 1 broker, got %d", len(l))
	}
	bRet, err := s.GetBroker(BROKER_UUID)
	if err != nil {
		t.Fatalf("GetBroker failed: %#v", err)
	}
	if bRet == nil {
		t.Fatalf("Did not get back a broker")
	}
	if strings.Compare(bRet.Name, b.Name) != 0 {
		t.Fatalf("Names don't match, expected: '%s', got '%s'", b.Name, bRet.Name)
	}
	catRet, err := s.GetInventory()
	if err != nil {
		t.Fatalf("Failed to get inventory: %#v", err)
	}
	if len(catRet.Services) != 2 {
		t.Fatalf("Expected 2 services from GetInventory, got %s ", len(catRet.Services))
	}

	err = s.DeleteBroker(BROKER_UUID)
	if err != nil {
		t.Fatalf("Failed to delete broker: %s : %#v", BROKER_UUID, err)
	}
	l, err = s.ListBrokers()
	if len(l) != 1 {
		t.Fatalf("Expected 1 broker, got %d", len(l))
	}
	bRet, err = s.GetBroker(BROKER_UUID)
	if err == nil {
		t.Fatalf("GetBroker returned a broker when there should be none")
	}
	bRet, err = s.GetBroker(BROKER_UUID_2)
	if err != nil {
		t.Fatalf("GetBroker failed for entry that should be there")
	}

	if bRet == nil {
		t.Fatalf("Did not get back a broker")
	}
	if strings.Compare(bRet.Name, b2.Name) != 0 {
		t.Fatalf("Names don't match, expected: '%s', got '%s'", b2.Name, bRet.Name)
	}
	catRet, err = s.GetInventory()
	if err != nil {
		t.Fatalf("Failed to get inventory: %#v", err)
	}
	if len(catRet.Services) != 1 {
		t.Fatalf("Expected 1 service from GetInventory, got %s ", len(catRet.Services))
	}
}
