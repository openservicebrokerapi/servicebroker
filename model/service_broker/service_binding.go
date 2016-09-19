package model

type CreateServiceBindingRequest struct {
	ServiceID    string                 `json:"service_id,omitempty"`
	PlanID       string                 `json:"plan_id,omitempty"`
	AppGUID      string                 `json:"app_guid,omitempty"`
	BindResource map[string]interface{} `json:"bind_resource,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

type CreateServiceBindingResponse struct {
	Credentials     interface{}   `json:"credentials"`
	SyslogDrainURL  string        `json:"syslog_drain_url, omitempty"`
	RouteServiceURL string        `json:"route_service_url, omitempty"`
	VolumeMounts    []interface{} `json:"volume_mounts,omitempty"`
}

type Credential struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
