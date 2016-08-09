package mem

import (
	"fmt"

	"github.com/cncf/servicebroker/k8s/service_controller/model"
	"github.com/cncf/servicebroker/k8s/service_controller/server"
)

type BindingPair struct {
	Binding    *model.ServiceBinding
	Credential *model.Credential
}

type InMemServiceStorage struct {
	brokers map[string]*model.ServiceBroker
	// This gets fetched when a SB is created (or possibly later when refetched).
	// It's static for now to keep compatibility, seems like this could be more dynamic.
	catalogs map[string]*model.Catalog
	// maps instance ID to instance
	services map[string]*model.ServiceInstance
	// maps binding ID to binding
	// TODO: support looking up all bindings for a service instance.
	bindings map[string]*BindingPair
}

var _ server.ServiceStorage = (*InMemServiceStorage)(nil)

func CreateInMemServiceStorage() server.ServiceStorage {
	return &InMemServiceStorage{
		brokers:  make(map[string]*model.ServiceBroker),
		catalogs: make(map[string]*model.Catalog),
		services: make(map[string]*model.ServiceInstance),
		bindings: make(map[string]*BindingPair),
	}
}

func (s *InMemServiceStorage) GetInventory() (*model.Catalog, error) {
	services := []*model.Service{}
	for _, v := range s.catalogs {
		services = append(services, v.Services...)
	}
	return &model.Catalog{Services: services}, nil
}

func (s *InMemServiceStorage) ListBrokers() ([]*model.ServiceBroker, error) {
	b := []*model.ServiceBroker{}
	for _, v := range s.brokers {
		b = append(b, v)
	}
	return b, nil
}

func (s *InMemServiceStorage) GetBroker(id string) (*model.ServiceBroker, error) {
	if b, ok := s.brokers[id]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("No such broker: %s", id)
}

func (s *InMemServiceStorage) GetBrokerByService(id string) (*model.ServiceBroker, error) {
	for k, v := range s.catalogs {
		for _, service := range v.Services {
			if service.ID == id {
				return s.brokers[k], nil
			}
		}
	}

	return nil, fmt.Errorf("No service matching ID %s", id)
}

func (s *InMemServiceStorage) AddBroker(broker *model.ServiceBroker, catalog *model.Catalog) error {
	if _, ok := s.brokers[broker.GUID]; ok {
		return fmt.Errorf("Broker %s already exists", broker.Name)
	}
	s.brokers[broker.GUID] = broker
	s.catalogs[broker.GUID] = catalog
	return nil
}

func (s *InMemServiceStorage) DeleteBroker(string) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *InMemServiceStorage) ServiceExists(id string) bool {
	_, err := s.GetService(id)
	return err == nil
}

func (s *InMemServiceStorage) ListServices() ([]*model.ServiceInstance, error) {
	services := []*model.ServiceInstance{}
	for _, v := range s.services {
		services = append(services, v)
	}
	return services, nil
}

func (s *InMemServiceStorage) GetService(id string) (*model.ServiceInstance, error) {
	service, ok := s.services[id]
	if !ok {
		return &model.ServiceInstance{}, fmt.Errorf("Service %s does not exist", id)
	}

	return service, nil
}

func (s *InMemServiceStorage) AddService(si *model.ServiceInstance) error {
	if s.ServiceExists(si.ID) {
		return fmt.Errorf("Service %s already exists", si.ID)
	}

	s.services[si.ID] = si
	return nil
}

func (s *InMemServiceStorage) DeleteService(id string) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *InMemServiceStorage) ListServiceBindings() ([]*model.ServiceBinding, error) {
	return []*model.ServiceBinding{}, fmt.Errorf("Not implemented yet")
}

func (s *InMemServiceStorage) GetServiceBinding(id string) (*model.Credential, error) {
	b, ok := s.bindings[id]
	if !ok {
		// TODO: this should return binding rather than credential.
		return &model.Credential{}, fmt.Errorf("Binding %s does not exist", id)
	}

	return b.Credential, nil
}

func (s *InMemServiceStorage) AddServiceBinding(binding *model.ServiceBinding, cred *model.Credential) error {
	_, err := s.GetServiceBinding(binding.ID)
	if err == nil {
		return fmt.Errorf("Binding %s already exists", binding.ID)
	}

	s.bindings[binding.ID] = &BindingPair{Binding: binding, Credential: cred}
	return nil
}

func (s *InMemServiceStorage) DeleteServiceBinding(id string) error {
	return fmt.Errorf("Not implemented yet")
}
