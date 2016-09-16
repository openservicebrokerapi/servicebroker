package model

// The Broker interface provides functions to deal with brokers.
type Broker interface {
	ListBrokers() ([]string, error)
	AddBroker(*ServiceBroker) error
	GetBroker(string) (*ServiceBroker, error)
	SetBroker(*ServiceBroker) error
	DeleteBroker(string) error
}

// Servicer provides functions to store services
type Servicer interface {
	ListServices() ([]string, error)
	AddService(*Service) error
	GetService(string) (*Service, error)
	SetService(*Service) error
	DeleteService(string) error
}

// Planner provides functions to store services
type Planner interface {
	ListPlans() ([]string, error)
	AddPlan(*ServicePlan) error
	GetPlan(string) (*ServicePlan, error)
	SetPlan(*ServicePlan) error
	DeletePlan(string) error
}

// The Instancer interface provides functions to deal with service instances.
type Instancer interface {
	ListInstances() ([]string, error)
	GetInstance(string) (*ServiceInstance, error)
	AddInstance(*ServiceInstance) error
	SetInstance(*ServiceInstance) error
	DeleteInstance(string) error
}

// The Binder interface provides functions to deal with service
// bindings.
type Binder interface {
	ListBindings() ([]string, error)
	GetBinding(string) (*ServiceBinding, error)
	AddBinding(*ServiceBinding) error
	SetBinding(*ServiceBinding) error
	DeleteBinding(string) error
}

// The ServiceStorage interface provides a comprehensive combined
// resource for end to end dealings with service brokers, service instances,
// and service bindings.
type ServiceStorage interface {
	Broker
	Servicer
	Planner
	Instancer
	Binder
}
