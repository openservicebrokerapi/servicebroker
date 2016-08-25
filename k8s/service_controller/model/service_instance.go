package model

type ServiceInstanceData struct {
	Instance *ServiceInstance

	// Bindings to pass to broker when instantiating the instance. If this is
	// not set at instance creation time, no bindings will ever be passed.
	//
	// Map of service name being bound to to credentials.
	Bindings map[string]*Credential
}

type ServiceInstance struct {
	Name             string `json:"name"`
	ID               string `json:"id"`
	DashboardURL     string `json:"dashboard_url"`
	InternalID       string `json:"internal_id, omitempty"`
	ServiceID        string `json:"service_id"`
	PlanID           string `json:"plan_id"`
	OrganizationGuid string `json:"organization_guid"`
	SpaceGuid        string `json:"space_guid"`

	LastOperation *LastOperation `json:"last_operation, omitempty"`

	Parameters interface{} `json:"parameters, omitempty"`
}

type LastOperation struct {
	State                    string `json:"state"`
	Description              string `json:"description"`
	AsyncPollIntervalSeconds int    `json:"async_poll_interval_seconds, omitempty"`
}

type CreateServiceInstanceRequest struct {
	Name              string                 `json:"name"`
	OrgID             string                 `json:"organization_guid"`
	ServicePlanGUID   string                 `json:"service_plan_guid"`
	SpaceID           string                 `json:"space_guid"`
	Parameters        map[string]interface{} `json:"parameters"`
	AcceptsIncomplete bool                   `json:"accepts_incomplete"`
}

type CreateServiceInstanceResponse struct {
	DashboardURL  string         `json:"dashboard_url, omitempty"`
	LastOperation *LastOperation `json:"last_operation, omitempty"`
}
