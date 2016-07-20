package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cncf/servicebroker/k8s/service_controller/model"
	"github.com/cncf/servicebroker/k8s/service_controller/utils"
)

const (
	CATALOG_URL_FMT_STR             = "%s/v2/catalog"
	CREATE_SERVICE_INSTANCE_FMT_STR = "%s/v2/service_instances/%s"
	BIND_FMT_STR                    = "%s/v2/service_instances/%s/service_bindings/%s"
)

type Controller struct {
	serviceStorage ServiceStorage
}

func CreateController(serviceStorage ServiceStorage) *Controller {
	return &Controller{
		serviceStorage: serviceStorage,
	}
}

func (c *Controller) ListServiceBrokers(w http.ResponseWriter, r *http.Request) {
	l, err := c.serviceStorage.ListBrokers()
	if err != nil {
		utils.WriteResponse(w, 400, err)
		return
	}
	utils.WriteResponse(w, 200, l)
}

func (c *Controller) GetServiceBroker(w http.ResponseWriter, r *http.Request) {
	name := utils.ExtractVarFromRequest(r, "broker_name")
	fmt.Printf("GetServiceBroker: %s", name)

	b, err := c.serviceStorage.GetBroker(name)
	if err != nil {
		fmt.Printf("Got Error: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	utils.WriteResponse(w, 200, b)
}

func (c *Controller) Inventory(w http.ResponseWriter, r *http.Request) {
	name := utils.ExtractVarFromRequest(r, "broker_name")
	fmt.Printf("Inventory: %s", name)

	i, err := c.serviceStorage.GetInventory(name)
	if err != nil {
		fmt.Printf("Got Error: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	utils.WriteResponse(w, 200, i)
}

func (c *Controller) CreateServiceBroker(w http.ResponseWriter, r *http.Request) {
	var b model.ServiceBroker
	err := utils.BodyToObject(r, &b)
	if err != nil {
		fmt.Printf("Error unmarshaling: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// Fetch the catalog from the broker
	u := fmt.Sprintf(CATALOG_URL_FMT_STR, b.BrokerURL)
	req, err := http.NewRequest("GET", u, nil)
	req.SetBasicAuth(b.AuthUsername, b.AuthPassword)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to fetch catalog from %s", u)
		utils.WriteResponse(w, 400, err)
		return
	}

	var catalog model.Catalog
	err = utils.ResponseBodyToObject(resp, &catalog)
	if err != nil {
		fmt.Printf("Failed to unmarshal catalog: ", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	err = c.serviceStorage.AddBroker(&b, &catalog)
	utils.WriteResponse(w, 200, b)
}

func (c *Controller) DeleteServiceBroker(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 400, "IMPLEMENT ME")
}

func (c *Controller) ListServiceInstances(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 400, "IMPLEMENT ME")
}

func (c *Controller) GetServiceInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting Service Instance")
	brokerName := utils.ExtractVarFromRequest(r, "broker_name")
	serviceName := utils.ExtractVarFromRequest(r, "service_name")

	si, err := c.serviceStorage.GetService(brokerName, serviceName)
	if err != nil {
		fmt.Printf("Couldn't fetch the service: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	utils.WriteResponse(w, 200, si)
}

func (c *Controller) CreateServiceInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Creating Service Instance")
	brokerName := utils.ExtractVarFromRequest(r, "broker_name")
	serviceName := utils.ExtractVarFromRequest(r, "service_name")

	if c.serviceStorage.ServiceExists(brokerName, serviceName) {
		err := fmt.Errorf("Service %s:%s already exists", brokerName, serviceName)
		fmt.Printf("%#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// Grab the broker to make sure it exists...
	broker, err := c.serviceStorage.GetBroker(brokerName)
	if err != nil {
		fmt.Printf("Couldn't fetch the broker: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	var req CreateServiceInstanceRequest
	err = utils.BodyToObject(r, &req)
	if err != nil {
		fmt.Printf("Error unmarshaling: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	serviceId, planId, err := c.getServiceAndPlanIds(brokerName, serviceName, req.PlanName)
	fmt.Printf("Found %s/%s => %s/%s", serviceName, req.PlanName, serviceId, planId)

	// Then actually make the request to reify the service instance
	createReq := &ServiceInstanceRequest{
		ServiceId:  serviceId,
		PlanId:     planId,
		Parameters: req.Parameters,
	}

	jsonBytes, err := json.Marshal(createReq)
	if err != nil {
		fmt.Printf("Failed to marshal: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	url := fmt.Sprintf(CREATE_SERVICE_INSTANCE_FMT_STR, broker.BrokerURL, serviceId)

	// TODO: Handle the auth
	createHttpReq, err := http.NewRequest("PUT", url, bytes.NewReader(jsonBytes))
	client := &http.Client{}
	fmt.Printf("Doing a request to: %s\n", url)
	resp, err := client.Do(createHttpReq)
	if err != nil {
		fmt.Printf("Failed to PUT: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	defer resp.Body.Close()

	si := model.ServiceInstance{}
	err = utils.ResponseBodyToObject(resp, &si)
	// TODO: Fix response to actually contain serviceId.
	si.ServiceId = serviceId

	c.serviceStorage.AddService(broker.Name, &si)

	utils.WriteResponse(w, 200, si)
}

func (c *Controller) DeleteServiceInstance(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 400, "IMPLEMENT ME")
}

func (c *Controller) ListServiceBindings(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 400, "IMPLEMENT ME")
}

func (c *Controller) GetServiceBinding(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting Service Binding")
	brokerName := utils.ExtractVarFromRequest(r, "broker_name")
	serviceName := utils.ExtractVarFromRequest(r, "service_name")
	bindingId := utils.ExtractVarFromRequest(r, "service_binding_guid")

	b, err := c.serviceStorage.GetServiceBinding(brokerName, serviceName, bindingId)
	if err != nil {
		fmt.Printf("%#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	utils.WriteResponse(w, 400, b)
}

func (c *Controller) CreateServiceBinding(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Creating Service Binding")
	brokerName := utils.ExtractVarFromRequest(r, "broker_name")
	serviceName := utils.ExtractVarFromRequest(r, "service_name")
	bindingId := utils.ExtractVarFromRequest(r, "service_binding_guid")

	if !c.serviceStorage.ServiceExists(brokerName, serviceName) {
		err := fmt.Errorf("Service %s:%s does not exist", brokerName, serviceName)
		fmt.Printf("%#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// Grab the broker to make sure it exists...
	broker, err := c.serviceStorage.GetBroker(brokerName)
	if err != nil {
		fmt.Printf("Couldn't fetch the broker: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	var req CreateServiceInstanceRequest
	err = utils.BodyToObject(r, &req)
	if err != nil {
		fmt.Printf("Error unmarshaling: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	serviceId, planId, err := c.getServiceAndPlanIds(brokerName, serviceName, req.PlanName)
	fmt.Printf("Found %s/%s => %s/%s", serviceName, req.PlanName, serviceId, planId)

	// Then actually make the request to reify the service instance
	createReq := &BindingRequest{
		ServiceId:  serviceId,
		PlanId:     planId,
		Parameters: req.Parameters,
	}

	jsonBytes, err := json.Marshal(createReq)
	if err != nil {
		fmt.Printf("Failed to marshal: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	url := fmt.Sprintf(BIND_FMT_STR, broker.BrokerURL, serviceId, bindingId)

	// TODO: Handle the auth
	createHttpReq, err := http.NewRequest("PUT", url, bytes.NewReader(jsonBytes))
	client := &http.Client{}
	fmt.Printf("Doing a request to: %s\n", url)
	resp, err := client.Do(createHttpReq)
	if err != nil {
		fmt.Printf("Failed to PUT: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	defer resp.Body.Close()

	sbr := model.CreateServiceBindingResponse{}
	err = utils.ResponseBodyToObject(resp, &sbr)
	if err != nil {
		fmt.Printf("Failed to unmarshal: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// TODO: get broker to actually return these values as part of response.
	sb := model.ServiceBinding{Id: bindingId, ServiceId: serviceId}

	c.serviceStorage.AddServiceBinding(broker.Name, &sb, &sbr.Credentials)
	utils.WriteResponse(w, 200, sb)
}

func (c *Controller) DeleteServiceBinding(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 400, "IMPLEMENT ME")
}

func (c *Controller) getServiceAndPlanIds(brokerName string, serviceName string, planName string) (string, string, error) {
	i, err := c.serviceStorage.GetInventory(brokerName)

	if err != nil {
		return "", "", err
	}
	var serviceFound = true

	for _, s := range i.Services {
		if strings.Compare(serviceName, s.Name) == 0 {
			// Ok, this is the service the customer is asking for, see if we find the plan...
			for _, p := range s.Plans {
				if strings.Compare(planName, p.Name) == 0 {
					return s.Id, p.Id, nil
				}
			}
		}
	}
	if !serviceFound {
		return "", "", fmt.Errorf("No service with name: '%s' found", serviceName)
	} else {
		return "", "", fmt.Errorf("No plan with name: '%s' found", planName)
	}
}

// This is what we get sent to us
type CreateServiceInstanceRequest struct {
	Name       string                 `json:"name"`
	PlanName   string                 `json:"plan"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

type ServiceInstanceRequest struct {
	OrgId             string                 `json:"organization_guid,omitempty"`
	PlanId            string                 `json:"plan_id,omitempty"`
	ServiceId         string                 `json:"service_id,omitempty"`
	SpaceId           string                 `json:"space_id,omitempty"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
	AcceptsIncomplete bool                   `json:"accepts_incomplete,omitempty"`
}

type BindingRequest struct {
	AppGuid      string                 `json:"app_guid,omitempty"`
	PlanId       string                 `json:"plan_id,omitempty"`
	ServiceId    string                 `json:"service_id,omitempty"`
	BindResource map[string]interface{} `json:"bind_resource,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}
