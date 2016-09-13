package mem

import (
	"fmt"

	model "github.com/servicebroker/servicebroker/k8s/service_controller/model"
)

type InMemServiceStorage struct {
	brokers   map[string]*model.ServiceBroker // key=id
	services  map[string]*model.Service
	plans     map[string]*model.ServicePlan
	instances map[string]*model.ServiceInstance
	bindings  map[string]*model.ServiceBinding
}

var _ model.ServiceStorage = (*InMemServiceStorage)(nil)

func CreateInMemServiceStorage() model.ServiceStorage {
	return &InMemServiceStorage{
		brokers:   map[string]*model.ServiceBroker{},
		services:  map[string]*model.Service{},
		plans:     map[string]*model.ServicePlan{},
		instances: map[string]*model.ServiceInstance{},
		bindings:  map[string]*model.ServiceBinding{},
	}
}

/* BROKERS */
/***********/

func (s *InMemServiceStorage) ListBrokers() ([]string, error) {
	b := []string{}
	for _, v := range s.brokers {
		b = append(b, v.ID)
	}

	return b, nil
}

func (s *InMemServiceStorage) AddBroker(broker *model.ServiceBroker) error {
	if _, ok := s.brokers[broker.ID]; ok {
		return fmt.Errorf("Broker %q already exists", broker.Name)
	}
	s.brokers[broker.ID] = broker
	return nil
}

func (s *InMemServiceStorage) GetBroker(id string) (*model.ServiceBroker, error) {
	if b, ok := s.brokers[id]; ok {
		return b, nil
	}
	return nil, nil
}

func (s *InMemServiceStorage) SetBroker(broker *model.ServiceBroker) error {
	s.brokers[broker.ID] = broker
	return nil
}

func (s *InMemServiceStorage) DeleteBroker(id string) error {
	_, err := s.GetBroker(id)
	if err != nil {
		return fmt.Errorf("Broker %q does not exist", id)
	}

	delete(s.brokers, id)
	return nil
}

/* SERVICES */
/************/

func (s *InMemServiceStorage) ListServices() ([]string, error) {
	services := []string{}
	for _, s := range s.services {
		services = append(services, s.ID)
	}
	return services, nil
}

func (s *InMemServiceStorage) AddService(service *model.Service) error {
	if _, ok := s.services[service.ID]; ok {
		return fmt.Errorf("Service %s already exists", service.ID)
	}

	s.services[service.ID] = service
	return nil
}

func (s *InMemServiceStorage) GetService(id string) (*model.Service, error) {
	service, ok := s.services[id]
	if !ok {
		return nil, nil
	}
	return service, nil
}

func (s *InMemServiceStorage) SetService(service *model.Service) error {
	s.services[service.ID] = service
	return nil
}

func (s *InMemServiceStorage) DeleteService(id string) error {
	delete(s.services, id)
	return nil
}

/* PLANS */
/*********/

func (s *InMemServiceStorage) ListPlans() ([]string, error) {
	plans := []string{}
	for _, p := range s.plans {
		plans = append(plans, p.ID)
	}
	return plans, nil
}

func (s *InMemServiceStorage) AddPlan(plan *model.ServicePlan) error {
	if _, ok := s.plans[plan.ID]; ok {
		return fmt.Errorf("Plan %s already exists", plan.ID)
	}

	s.plans[plan.ID] = plan
	return nil
}

func (s *InMemServiceStorage) GetPlan(id string) (*model.ServicePlan, error) {
	plan, ok := s.plans[id]
	if !ok {
		return nil, nil
	}
	return plan, nil
}

func (s *InMemServiceStorage) SetPlan(plan *model.ServicePlan) error {
	s.plans[plan.ID] = plan
	return nil
}

func (s *InMemServiceStorage) DeletePlan(id string) error {
	delete(s.plans, id)
	return nil
}

/* INSTANCES */
/*************/

func (s *InMemServiceStorage) ListInstances() ([]string, error) {
	instances := []string{}
	for _, i := range s.instances {
		instances = append(instances, i.ID)
	}
	return instances, nil
}

func (s *InMemServiceStorage) AddInstance(instance *model.ServiceInstance) error {
	if _, ok := s.instances[instance.ID]; ok {
		return fmt.Errorf("Insance %q already exists", instance.ID)
	}

	s.instances[instance.ID] = instance
	return nil
}

func (s *InMemServiceStorage) GetInstance(id string) (*model.ServiceInstance, error) {
	instance, ok := s.instances[id]
	if !ok {
		return nil, nil
	}
	return instance, nil
}

func (s *InMemServiceStorage) SetInstance(instance *model.ServiceInstance) error {
	s.instances[instance.ID] = instance
	return nil
}

func (s *InMemServiceStorage) DeleteInstance(id string) error {
	delete(s.instances, id)
	return nil
}

/* BINDINGS */
/************/

func (s *InMemServiceStorage) ListBindings() ([]string, error) {
	bindings := []string{}
	for _, v := range s.bindings {
		bindings = append(bindings, v.ID)
	}
	return bindings, nil
}

func (s *InMemServiceStorage) GetBinding(id string) (*model.ServiceBinding, error) {
	b, ok := s.bindings[id]
	if !ok {
		return nil, nil
	}
	return b, nil
}

func (s *InMemServiceStorage) AddBinding(binding *model.ServiceBinding) error {
	_, err := s.GetBinding(binding.ID)
	if err != nil {
		return fmt.Errorf("Binding %q already exists", binding.ID)
	}

	s.bindings[binding.ID] = binding
	return nil
}

func (s *InMemServiceStorage) SetBinding(binding *model.ServiceBinding) error {
	s.bindings[binding.ID] = binding
	return nil
}

func (s *InMemServiceStorage) DeleteBinding(id string) error {
	delete(s.bindings, id)
	return nil
}
