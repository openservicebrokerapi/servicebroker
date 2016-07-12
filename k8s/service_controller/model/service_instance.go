package model

type ServiceInstance struct {
	Id               string `json:"id"`
	DashboardUrl     string `json:"dashboard_url"`
	InternalId       string `json:"internalId, omitempty"`
	ServiceId        string `json:"service_id"`
	PlanId           string `json:"plan_id"`
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

type CreateServiceInstanceResponse struct {
	DashboardUrl  string         `json:"dashboard_url, omitempty"`
	LastOperation *LastOperation `json:"last_operation, omitempty"`
}
