package model

type ServiceBinding struct {
	ID                      string                 `json:"id"`
	FromServiceInstanceName string                 `json:"from_service_instance_name"`
	ServiceInstanceGUID     string                 `json:"service_instance_guid"`
	Parameters              map[string]interface{} `json:"parameters,omitempty"`
}

type CreateServiceBindingResponse struct {
	// SyslogDrainURL string      `json:"syslog_drain_url, omitempty"`
	Credentials *interface{} `json:"credentials"`
}
