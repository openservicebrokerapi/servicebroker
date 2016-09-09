package model

type ServiceInstance struct {
	// Fields from the CLI
	ID                string
	Name              string
	PlanGUID          string
	SpaceGUID         string
	Parameters        map[string]interface{}
	Tags              []string
	AcceptsIncomplete bool

	// Fields from the SB
	DashboardURL  string
	LastOperation *LastOperation

	// Internals
	Bindings map[string]*interface{}
	Service  *Service
}

type LastOperation struct {
	State                    string `json:"state"`
	Description              string `json:"description"`
	AsyncPollIntervalSeconds int    `json:"async_poll_interval_seconds, omitempty"`
}
