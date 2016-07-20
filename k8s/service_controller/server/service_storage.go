package server

import (
	"github.com/cncf/servicebroker/k8s/service_controller/model"
)

// The Broker interface provides functions to deal with brokers.
type Broker interface {
	ListBrokers() ([]*model.ServiceBroker, error)
	GetBroker(string) (*model.ServiceBroker, error)
	GetInventory(string) (*model.Catalog, error)
	AddBroker(*model.ServiceBroker, *model.Catalog) error
	DeleteBroker(string) error
}

// The Instancer interface provides functions to deal with service instances.
type Instancer interface {
	ListServices(string) ([]*model.ServiceInstance, error)
	GetService(string, string) (*model.ServiceInstance, error)
	ServiceExists(string, string) bool
	AddService(string, *model.ServiceInstance) error
	DeleteService(string, string) error
}

// The Binder interface provides functions to deal with service
// bindings.
type Binder interface {
	ListServiceBindings(string, string) ([]*model.ServiceBinding, error)
	GetServiceBinding(string, string, string) (*model.Credential, error)
	AddServiceBinding(string, *model.ServiceBinding, *model.Credential) error
	DeleteServiceBinding(string, string, string) error
}

// The ServiceStorage interface provides a comprehensive combined
// resource for end to end dealings with service brokers, service instances,
// and service bindings.
type ServiceStorage interface {
	Broker
	Instancer
	Binder
}
