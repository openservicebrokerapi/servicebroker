package model

type ServiceBinding struct {
	// From Broker
	Credentials     interface{}
	SyslogDrainURL  string
	RouteServiceURL string
	VolumeMounts    []interface{}

	// From CLI
	AppName             string
	ServiceInstanceName string                 `json:"service_instance_name"`
	ServiceInstanceID   string                 `json:"service_instance_id"`
	Parameters          map[string]interface{} `json:"parameters,omitempty"`

	// Our extras
	ID string
}
