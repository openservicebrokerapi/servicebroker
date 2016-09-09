package model

type ServiceInstance struct {
	// From CF spec
	Name            string          `json:"name"`
	Credentials     string          `json:"credentials"`
	ServicePlanGUID string          `json:"service_plan_guid"`
	SpaceGUID       string          `json:"space_guid"`
	DashboardURL    string          `json:"dashboard_url"`
	Type            string          `json:"type"`
	LastOperation   *LastOperation1 `json:"last_operation, omitempty"`
	SpaceURL        string          `json:"space_url"`
	ServicePlanURL  string          `json:"service_plan_url"`
	RoutesURL       string          `json:"routes_url"`
	Tags            []string        `json:"tags"`

	Parameters interface{} `json:"parameters, omitempty"`

	// Our extras
	ID        string `json:"id"`
	ServiceID string `json:"service_id"`
}

type LastOperation1 struct {
	Type        string `json:"state"`
	State       string `json:"state"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
}
type CreateServiceInstanceRequest struct {
	Name            string                 `json:"name"`
	ServicePlanGUID string                 `json:"service_plan_guid"`
	SpaceID         string                 `json:"space_guid"`
	Parameters      map[string]interface{} `json:"parameters"`
	Tags            []string               `json:"tags"`
}

type CreateServiceInstanceResponse struct {
	DashboardURL  string          `json:"dashboard_url, omitempty"`
	LastOperation *LastOperation2 `json:"last_operation, omitempty"`
}

type LastOperation2 struct {
	State                    string `json:"state"`
	Description              string `json:"description"`
	AsyncPollIntervalSeconds int    `json:"async_poll_interval_seconds, omitempty"`
}
