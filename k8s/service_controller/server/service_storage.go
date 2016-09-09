package server

import (
	model "github.com/servicebroker/servicebroker/k8s/service_controller/model"
)

// The Broker interface provides functions to deal with brokers.
type Broker interface {
	ListBrokers() ([]*model.ServiceBroker, error)
	GetBroker(string) (*model.ServiceBroker, error)
	GetBrokerByService(string) (*model.ServiceBroker, error)
	GetInventory() (*model.Catalog, error)
	AddBroker(*model.ServiceBroker, *model.Catalog) error
	DeleteBroker(string) error
}

// The Instancer interface provides functions to deal with service instances.
type Instancer interface {
	ListServices() ([]*model.ServiceInstance, error)
	GetService(string) (*model.ServiceInstance, error)
	ServiceExists(string) bool
	AddService(*model.ServiceInstance) error
	SetService(*model.ServiceInstance) error
	DeleteService(string) error
}

// The Binder interface provides functions to deal with service
// bindings.
type Binder interface {
	ListServiceBindings() ([]*model.ServiceBinding, error)
	GetServiceBinding(string) (*model.ServiceBinding, error)
	AddServiceBinding(*model.ServiceBinding, *interface{}) error
	DeleteServiceBinding(string) error
}

// The ServiceStorage interface provides a comprehensive combined
// resource for end to end dealings with service brokers, service instances,
// and service bindings.
type ServiceStorage interface {
	Broker
	Instancer
	Binder
	// This provides access to the available services provided by
	// all known brokers. Equivalent to `cf marketplace`.
	GetServices() ([]*model.Service, error)
}
