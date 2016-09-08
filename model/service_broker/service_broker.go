package model

type CreateServiceBrokerRequest struct {
	Name         string `json:"name"`
	BrokerURL    string `json:"broker_url"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
	SpaceGUID    string `json:"space_guid"` // CF-specific - FIXME
}

type CreateServiceBrokerResponse struct {
	Metadata ServiceBrokerMetadata `json:"metadata"`
	Entity   ServiceBrokerEntity   `json:"entity"`
}

type ServiceBrokerMetadata struct {
	GUID      string `json:"guid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
	URL       string `json:"url"`
}

type ServiceBrokerEntity struct {
	Name         string `json:"name"`
	BrokerURL    string `json:"broker_url"`
	AuthUsername string `json:"auth_username"`
	// space_guid
}
