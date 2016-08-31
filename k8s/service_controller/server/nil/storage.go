package nil

import (
	"fmt"

	"github.com/servicebroker/servicebroker/k8s/service_controller/model"
	"github.com/servicebroker/servicebroker/k8s/service_controller/server"
)

// when you absolutely have to have storage, but you don't really care what happens.
type NilServiceStorage struct {
}

func CreateNilServiceStorage() server.ServiceStorage {
	return &NilServiceStorage{}
}

func (s *NilServiceStorage) ListBrokers() ([]*model.ServiceBroker, error) {
	return nil, nil
}

func (s *NilServiceStorage) GetBroker(id string) (*model.ServiceBroker, error) {
	return nil, nil
}

func (s *NilServiceStorage) GetInventory() (*model.Catalog, error) {
	return nil, nil
}

func (s *NilServiceStorage) AddBroker(broker *model.ServiceBroker, catalog *model.Catalog) error {
	return nil
}

func (s *NilServiceStorage) DeleteBroker(id string) error {
	return nil
}

func (s *NilServiceStorage) ServiceExists(id string) bool {
	return false
}

func (s *NilServiceStorage) ListServices() ([]*model.ServiceInstanceData, error) {
	return nil, nil
}

func (s *NilServiceStorage) GetService(id string) (*model.ServiceInstanceData, error) {
	return nil, nil
}

func (s *NilServiceStorage) AddService(si *model.ServiceInstanceData) error {
	return nil
}

func (s *NilServiceStorage) SetService(si *model.ServiceInstanceData) error {
	return nil
}

func (s *NilServiceStorage) DeleteService(id string) error {
	return nil
}

func (s *NilServiceStorage) ListServiceBindings() ([]*model.ServiceBinding, error) {
	return nil, nil
}

func (s *NilServiceStorage) GetServiceBinding(id string) (*model.ServiceBinding, error) {
	return nil, nil
}

func (s *NilServiceStorage) AddServiceBinding(binding *model.ServiceBinding, cred *model.Credential) error {
	return nil
}

func (s *NilServiceStorage) DeleteServiceBinding(id string) error {
	return nil
}

func (s *NilServiceStorage) GetBrokerByService(id string) (*model.ServiceBroker, error) {
	return nil, nil
}

// NotImplementedYetServiceStorage behaves appropriately, in that it always returns an error if possible.
type NotImplementedYetServiceStorage struct {
}

func CreateNotImplementedServiceStorage() server.ServiceStorage {
	return &NotImplementedYetServiceStorage{}
}

func (s *NotImplementedYetServiceStorage) ListBrokers() ([]*model.ServiceBroker, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) GetBroker(id string) (*model.ServiceBroker, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) GetBrokerByService(id string) (*model.ServiceBroker, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) GetInventory() (*model.Catalog, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) AddBroker(broker *model.ServiceBroker, catalog *model.Catalog) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) DeleteBroker(id string) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) ServiceExists(id string) bool {
	return false
}

func (s *NotImplementedYetServiceStorage) ListServices() ([]*model.ServiceInstance, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) GetServices() ([]*model.Service, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) GetService(id string) (*model.ServiceInstance, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) AddService(si *model.ServiceInstance) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) DeleteService(id string) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) ListServiceBindings() ([]*model.ServiceBinding, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) GetServiceBinding(id string) (*model.Credential, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) AddServiceBinding(binding *model.ServiceBinding, cred *model.Credential) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *NotImplementedYetServiceStorage) DeleteServiceBinding(id string) error {
	return fmt.Errorf("Not implemented yet")
}

var _ server.ServiceStorage = (*NotImplementedYetServiceStorage)(nil)
var _ server.ServiceStorage = (*NilServiceStorage)(nil)
