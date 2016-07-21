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
	brokerMap map[string]*model.ServiceBroker
	// This gets fetched when a SB is created (or possibly later when refetched).
	// It's static for now to keep compatibility, seems like this could be more dynamic.
	catalogs map[string]*model.Catalog
	// maps brokers to array of service instances
	serviceMap map[string][]*model.ServiceInstance
	// hacky 2 dimensional map, key is broker:service, returns array of service bindings
	bindingMap map[string][]*BindingPair
}

var _ server.ServiceStorage = (*InMemServiceStorage)(nil)

func CreateInMemServiceStorage() server.ServiceStorage {
	return &InMemServiceStorage{
		brokerMap:  make(map[string]*model.ServiceBroker),
		catalogs:   make(map[string]*model.Catalog),
		serviceMap: make(map[string][]*model.ServiceInstance),
		bindingMap: make(map[string][]*BindingPair),
	}
}

func (s *InMemServiceStorage) ListBrokers() ([]*model.ServiceBroker, error) {
	b := []*model.ServiceBroker{}
	for _, v := range s.brokerMap {
		b = append(b, v)
	}
	return b, nil
}
func (s *InMemServiceStorage) GetBroker(name string) (*model.ServiceBroker, error) {
	if b, ok := s.brokerMap[name]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("No such broker: %s", name)
}

func (s *InMemServiceStorage) GetInventory(name string) (*model.Catalog, error) {
	if b, ok := s.catalogs[name]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("No catalog for broker: %s", name)
}

func (s *InMemServiceStorage) AddBroker(broker *model.ServiceBroker, catalog *model.Catalog) error {
	if _, ok := s.brokerMap[broker.Name]; ok {
		return fmt.Errorf("Broker %s already exists", broker.Name)
	}
	s.brokerMap[broker.Name] = broker
	s.catalogs[broker.Name] = catalog
	s.serviceMap[broker.Name] = []*model.ServiceInstance{}
	return nil
}

func (s *InMemServiceStorage) DeleteBroker(string) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *InMemServiceStorage) ServiceExists(broker, service string) bool {
	_, err := s.GetService(broker, service)
	return err == nil
}

func (s *InMemServiceStorage) ListServices(broker string) ([]*model.ServiceInstance, error) {
	sm, ok := s.serviceMap[broker]
	if !ok {
		return []*model.ServiceInstance{}, fmt.Errorf("Broker %s does not exist", broker)
	}
	return sm, nil
}

func (s *InMemServiceStorage) GetService(broker, service string) (*model.ServiceInstance, error) {
	sm, ok := s.serviceMap[broker]
	if !ok {
		return &model.ServiceInstance{}, fmt.Errorf("Broker %s does not exist", broker)
	}

	for _, si := range sm {
		if si.ServiceId == service {
			return si, nil
		}
	}
	return &model.ServiceInstance{}, fmt.Errorf("Service %s not found in broker %s", service, broker)
}

func (s *InMemServiceStorage) AddService(broker string, si *model.ServiceInstance) error {
	if s.ServiceExists(broker, si.ServiceId) {
		return fmt.Errorf("Service %s already exists in broker %s", si.ServiceId, broker)
	}

	s.serviceMap[broker] = append(s.serviceMap[broker], si)
	return nil
}

func (s *InMemServiceStorage) DeleteService(string, string) error {
	return fmt.Errorf("Not implemented yet")
}

func (s *InMemServiceStorage) ListServiceBindings(string, string) ([]*model.ServiceBinding, error) {
	return []*model.ServiceBinding{}, fmt.Errorf("Not implemented yet")
}

func (s *InMemServiceStorage) GetServiceBinding(broker, service, binding string) (*model.Credential, error) {
	key := bindingKey(broker, service)
	bs, ok := s.bindingMap[key]
	if !ok {
		return &model.Credential{}, fmt.Errorf("Service %s in broker %s does not exist", service, broker)
	}

	for _, bp := range bs {
		if bp.Binding.Id == binding {
			return bp.Credential, nil
		}
	}
	return &model.Credential{}, fmt.Errorf("Binding %s not found in service %s", binding, service)
}

func (s *InMemServiceStorage) AddServiceBinding(broker string, binding *model.ServiceBinding, cred *model.Credential) error {
	_, err := s.GetServiceBinding(broker, binding.ServiceId, binding.Id)
	if err == nil {
		return fmt.Errorf("Binding %s already exists for service %s in broker %s", binding.Id, binding.ServiceId, broker)
	}

	key := bindingKey(broker, binding.ServiceId)
	s.bindingMap[key] = append(s.bindingMap[key], &BindingPair{Binding: binding, Credential: cred})
	return nil
}

func (s *InMemServiceStorage) DeleteServiceBinding(string, string, string) error {
	return fmt.Errorf("Not implemented yet")
}

func bindingKey(broker, service string) string {
	return fmt.Sprintf("%s:%s", broker, service)
}
