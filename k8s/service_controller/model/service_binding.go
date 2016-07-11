package model

type ServiceBinding struct {
	Id                string `json:"id"`
	ServiceId         string `json:"service_id"`
	AppId             string `json:"app_id"`
	ServicePlanId     string `json:"service_plan_id"`
	PrivateKey        string `json:"private_key"`
	ServiceInstanceId string `json:"service_instance_id"`
}

type CreateServiceBindingResponse struct {
	// SyslogDrainUrl string      `json:"syslog_drain_url, omitempty"`
	Credentials Credential `json:"credentials"`
}

type Credential struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
