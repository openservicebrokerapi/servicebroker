package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	scmodel "github.com/servicebroker/servicebroker/model/service_controller"
)

// TODO(vaikas): Move these into a ../lib or ../../lib?

func fetchInventory() (*scmodel.GetCatalogResponse, error) {
	u := fmt.Sprintf("%s%s", controller, SERVICE_PLANS_URL)
	i, err := callHttp(u, "GET", "inventory", nil)
	if err != nil {
		return nil, err
	}
	var c scmodel.GetCatalogResponse
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
	var brokers []scmodel.ServiceBroker
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

func fetchPrettyBindings() (string, error) {
	u := fmt.Sprintf("%s%s", controller, SERVICE_BINDINGS_URL)
	i, err := callHttp(u, "GET", "list service bindings", nil)
	if err != nil {
		return "", err
	}
	var bindings []scmodel.ServiceBinding
	err = json.Unmarshal([]byte(i), &bindings)
	if err != nil {
		return "", err
	}
	var ret []string

	for _, sb := range bindings {
		bar := fmt.Sprintf(SERVICE_INSTANCES_FMT_STR, sb.ServiceInstanceGUID)
		u = fmt.Sprintf("%s%s", controller, bar)
		i, err = callHttp(u, "GET", "describe service instance", nil)
		if err != nil {
			return "", err
		}
		var si scmodel.ServiceInstance
		err = json.Unmarshal([]byte(i), &si)
		if err != nil {
			return "", err
		}

		ret = append(ret, fmt.Sprintf("%s -> %s\n\t%+v\n", sb.FromServiceInstanceName, si.Name, sb.Parameters))
	}
	return strings.Join(ret, "\n"), nil
}

// Fetches the inventory from the SC and maps service:plan to the unique ID of the service plan
// This could be more efficient with client side caching, etc. but for now will suffice.
func fetchServicePlanGUID(service string, plan string) (string, error) {
	u := fmt.Sprintf("%s%s", controller, SERVICE_PLANS_URL)
	i, err := callHttp(u, "GET", "inventory", nil)
	if err != nil {
		return "", err
	}
	var c scmodel.GetCatalogResponse
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

// Fetches the GUID for a given serviceInstance (name) from the SC
// This could be more efficient with client side caching, etc. but for now will suffice.
func fetchServiceInstanceGUID(serviceInstance string) (string, error) {
	u := fmt.Sprintf("%s%s", controller, SERVICE_INSTANCES_URL)
	s, err := callHttp(u, "GET", "list service instances", nil)
	if err != nil {
		return "", err
	}
	var services []*scmodel.ServiceInstance
	err = json.Unmarshal([]byte(s), &services)
	if err != nil {
		return "", err
	}
	for _, service := range services {
		fmt.Printf("Checking: %#v\n", service)
		if service.Name == serviceInstance {
			return service.ID, nil
		}
	}
	return "", fmt.Errorf("Can't find a ServiceInstance : %s", serviceInstance)
}
