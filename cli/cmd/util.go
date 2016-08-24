package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cncf/servicebroker/k8s/service_controller/model"
)

// TODO(vaikas): Move these into a ../lib or ../../lib?

func fetchInventory() (*model.Catalog, error) {
	u := fmt.Sprintf("%s%s", controller, SERVICE_PLANS_URL)
	i, err := callHttp(u, "GET", "inventory", nil)
	if err != nil {
		return nil, err
	}
	var c model.Catalog
	err = json.Unmarshal([]byte(i), &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Maps a Broker name to UUID
func fetchBrokerGUID(broker string) (string, error) {
	u := fmt.Sprintf("%s%s", controller, BROKERS_URL)
	b, err := callHttp(u, "GET", "list brokers", nil)

	if err != nil {
		return "", err
	}
	var brokers []model.ServiceBroker
	err = json.Unmarshal([]byte(b), &brokers)
	if err != nil {
		return "", err
	}
	for _, s := range brokers {
		if s.Name == broker {
			return s.GUID, nil
		}
	}
	return "", fmt.Errorf("Can't find a broker : %s", broker)
}

// Fetches the inventory from the SC and maps service:plan to the unique ID of the service plan
// This could be more efficient with client side caching, etc. but for now will suffice.
func fetchServicePlanGUID(service string, plan string) (string, error) {
	u := fmt.Sprintf("%s%s", controller, SERVICE_PLANS_URL)
	i, err := callHttp(u, "GET", "inventory", nil)
	if err != nil {
		return "", err
	}
	var c model.Catalog
	err = json.Unmarshal([]byte(i), &c)
	if err != nil {
		return "", err
	}
	for _, s := range c.Services {
		if s.Name == service {
			for _, p := range s.Plans {
				if p.Name == plan {
					fmt.Printf("Found Service Plan GUID as %s for %s : %s", p.ID, service, plan)
					return p.ID, nil
				}
			}
		}
	}
	return "", fmt.Errorf("Can't find a service / plan : %s/%s", service, plan)
}
