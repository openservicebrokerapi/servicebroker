package model

type ServiceBinding struct {
	ID                string `json:"id"`
	ServiceID         string `json:"service_id"`
	AppID             string `json:"app_id"`
	ServicePlanID     string `json:"service_plan_id"`
	PrivateKey        string `json:"private_key"`
	ServiceInstanceID string `json:"service_instance_id"`
}

type CreateServiceBindingResponse struct {
	// SyslogDrainURL string      `json:"syslog_drain_url, omitempty"`
	Credentials Credential `json:"credentials"`
}

type Credential struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
