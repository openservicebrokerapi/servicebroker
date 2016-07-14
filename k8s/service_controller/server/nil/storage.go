package nil

import (
	"fmt"

	"github.com/cncf/servicebroker/k8s/service_controller/model"
	"github.com/cncf/servicebroker/k8s/service_controller/server"
)

// when you absolutely have to have storage, but you don't really care what happens.
type NilServiceStorage struct {
}

func CreateNilServiceStorage() server.ServiceStorage {
	return &NilServiceStorage{}
}

func (s *NilServiceStorage) ListBrokers() ([]*server.ServiceBroker, error) {
	return nil, nil
}

func (s *NilServiceStorage) GetBroker(name string) (*server.ServiceBroker, error) {
	return nil, nil
}

func (s *NilServiceStorage) GetInventory(name string) (*model.Catalog, error) {
	return nil, nil
}

func (s *NilServiceStorage) AddBroker(broker *server.ServiceBroker, catalog *model.Catalog) error {
	return nil
}

func (s *NilServiceStorage) DeleteBroker(string) error {
	return nil
}

func (s *NilServiceStorage) ServiceExists(broker, service string) bool {
	return false
}

func (s *NilServiceStorage) ListServices(broker string) ([]*model.ServiceInstance, error) {
	return nil, nil
}

func (s *NilServiceStorage) GetService(broker, service string) (*model.ServiceInstance, error) {
	return nil, nil
}

func (s *NilServiceStorage) AddService(broker string, si *model.ServiceInstance) error {
	return nil
}

func (s *NilServiceStorage) DeleteService(string, string) error {
	return nil
}

func (s *NilServiceStorage) ListServiceBindings(string, string) ([]*model.ServiceBinding, error) {
	return nil, nil
}

func (s *NilServiceStorage) GetServiceBinding(broker, service, binding string) (*model.Credential, error) {
	return nil, nil
}

func (s *NilServiceStorage) AddServiceBinding(broker string, binding *model.ServiceBinding, cred *model.Credential) error {
	return nil
}

func (s *NilServiceStorage) DeleteServiceBinding(string, string, string) error {
	return nil
}

// NotImplementedYetServiceStorage behaves appropriately, in that it always returns an error if possible.
type NotImplementedYetServiceStorage struct {
}

func CreateNotImplementedServiceStorage() server.ServiceStorage {
	return &NotImplementedYetServiceStorage{}
}

func (s *NotImplementedYetServiceStorage) ListBrokers() ([]*server.ServiceBroker, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) GetBroker(name string) (*server.ServiceBroker, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) GetInventory(name string) (*model.Catalog, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) AddBroker(broker *server.ServiceBroker, catalog *model.Catalog) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) DeleteBroker(string) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) ServiceExists(broker, service string) bool {
	return false
}

func (s *NotImplementedYetServiceStorage) ListServices(broker string) ([]*model.ServiceInstance, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) GetService(broker, service string) (*model.ServiceInstance, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) AddService(broker string, si *model.ServiceInstance) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) DeleteService(string, string) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) ListServiceBindings(string, string) ([]*model.ServiceBinding, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) GetServiceBinding(broker, service, binding string) (*model.Credential, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) AddServiceBinding(broker string, binding *model.ServiceBinding, cred *model.Credential) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) DeleteServiceBinding(string, string, string) error {
	return fmt.Errorf("Not implemented yet")
}

var _ server.ServiceStorage = (*NotImplementedYetServiceStorage)(nil)
var _ server.ServiceStorage = (*NilServiceStorage)(nil)
