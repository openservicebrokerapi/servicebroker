package model

type ServiceInstance struct {
	// Fields from the CLI
	ID                string
	Name              string
	Plan              string
	SpaceID           string
	Parameters        map[string]interface{}
	Tags              []string
	AcceptsIncomplete bool

	// Fields from the SB
	DashboardURL  string
	LastOperation LastOperation

	// Internals
	Service  string            // Service ID
	Bindings map[string]string // key=bindingID, binding ID - switch to {}
}

type LastOperation struct {
	UpdatedAt                string
	State                    string
	Description              string
	AsyncPollIntervalSeconds int
}
