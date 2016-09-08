package model

// http://apidocs.cloudfoundry.org/239/service_brokers/create_a_service_broker.html

type ServiceBroker struct {
	GUID         string
	Name         string
	BrokerURL    string
	AuthUsername string
	AuthPassword string
	// SpaceGUID    string

	Created int64 `json:",string"`
	Updated int64
	SelfURL string
}
