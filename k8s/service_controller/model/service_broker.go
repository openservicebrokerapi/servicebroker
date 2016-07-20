package model

type ServiceBroker struct {
	// json info from create docs
	// http://apidocs.cloudfoundry.org/239/service_brokers/create_a_service_broker.html
	Name         string `json:"name"`
	BrokerURL    string `json:"broker_url"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
	SpaceGUID    string `json:"space_guid"` // CF-specific - FIXME
}
