package model

type CreateServiceInstanceRequest struct {
	OrgID             string                 `json:"organization_guid,omitempty"`
	PlanID            string                 `json:"plan_id,omitempty"`
	ServiceID         string                 `json:"service_id,omitempty"`
	SpaceID           string                 `json:"space_guid,omitempty"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
	AcceptsIncomplete bool                   `json:"accepts_incomplete,omitempty"`
}

type CreateServiceInstanceResponse struct {
	DashboardURL  string         `json:"dashboard_url, omitempty"`
	LastOperation *LastOperation `json:"last_operation, omitempty"`
}

type GetLastOperationResponse struct {
	State       string `json:"state"`
	Description string `json:"description"`
}

type LastOperation struct {
	State       string `json:"state"`
	Description string `json:"description"`
}
