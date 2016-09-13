package model

type Service struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Bindable       bool     `json:"bindable"`
	PlanUpdateable bool     `json:"plan_updateable, omitempty"`
	Tags           []string `json:"tags, omitempty"`
	Requires       []string `json:"requires, omitempty"`

	Metadata        interface{} `json:"metadata, omitempty"`
	DashboardClient interface{} `json:"dashboard_client"`

	ServiceBroker string            // ServiceBroker ID
	Plans         []string          `json:"plans"` // ServicePlan ID
	Instances     map[string]string // ServiceIntance ID -- TODO make {}
}
