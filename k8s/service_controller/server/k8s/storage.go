package k8s

import (
	"fmt"

	"github.com/cncf/servicebroker/k8s/service_controller/model"
	"github.com/cncf/servicebroker/k8s/service_controller/server"
)

type K8sServiceStorage struct {
}

var _ server.ServiceStorage = (*K8sServiceStorage)(nil)

func CreateServiceStorage() server.ServiceStorage {
	return &K8sServiceStorage{}
}

func (s *K8sServiceStorage) ListBrokers() ([]*server.ServiceBroker, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) GetBroker(name string) (*server.ServiceBroker, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) GetInventory(name string) (*model.Catalog, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) AddBroker(broker *server.ServiceBroker, catalog *model.Catalog) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) DeleteBroker(string) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) ServiceExists(broker, service string) bool {
	return false
}

func (s *K8sServiceStorage) ListServices(broker string) ([]*model.ServiceInstance, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) GetService(broker, service string) (*model.ServiceInstance, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) AddService(broker string, si *model.ServiceInstance) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) DeleteService(string, string) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) ListServiceBindings(string, string) ([]*model.ServiceBinding, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) GetServiceBinding(broker, service, binding string) (*model.Credential, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) AddServiceBinding(broker string, binding *model.ServiceBinding, cred *model.Credential) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *K8sServiceStorage) DeleteServiceBinding(string, string, string) error {
	return fmt.Errorf("Not implemented yet")
}
