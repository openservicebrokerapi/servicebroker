package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cncf/servicebroker/k8s/service_controller/model"
	"github.com/cncf/servicebroker/k8s/service_controller/utils"
	"github.com/satori/go.uuid"
)

const (
	CATALOG_URL_FMT_STR             = "%s/v2/catalog"
	CREATE_SERVICE_INSTANCE_FMT_STR = "%s/v2/service_instances/%s"
	BIND_FMT_STR                    = "%s/v2/service_instances/%s/service_bindings/%s"
)

type Controller struct {
	storage ServiceStorage
}

func CreateController(storage ServiceStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}

//
// Inventory.
//

func (c *Controller) Inventory(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Inventory\n")

	i, err := c.storage.GetInventory()
	if err != nil {
		fmt.Printf("Got Error: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	utils.WriteResponse(w, 200, i)
}

//
// Service Broker.
//

func (c *Controller) ListServiceBrokers(w http.ResponseWriter, r *http.Request) {
	l, err := c.storage.ListBrokers()
	if err != nil {
		fmt.Printf("Got Error: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	utils.WriteResponse(w, 200, l)
}

func (c *Controller) GetServiceBroker(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractVarFromRequest(r, "broker_id")
	fmt.Printf("GetServiceBroker: %s\n", id)

	b, err := c.storage.GetBroker(id)
	if err != nil {
		fmt.Printf("Got Error: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	utils.WriteResponse(w, 200, b)
}

func (c *Controller) CreateServiceBroker(w http.ResponseWriter, r *http.Request) {
	var sbReq model.CreateServiceBrokerRequest
	err := utils.BodyToObject(r, &sbReq)
	if err != nil {
		fmt.Printf("Error unmarshaling: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	sb := model.ServiceBroker{
		GUID:         uuid.NewV4().String(),
		Name:         sbReq.Name,
		BrokerURL:    sbReq.BrokerURL,
		AuthUsername: sbReq.AuthUsername,
		AuthPassword: sbReq.AuthPassword,

		Created: time.Now().Unix(),
		Updated: 0,
		// SelfURL: "/v2/service_brokers/" + sb.GUID,
	}
	sb.SelfURL = "/v2/service_brokers/" + sb.GUID

	// Fetch the catalog from the broker
	u := fmt.Sprintf(CATALOG_URL_FMT_STR, sb.BrokerURL)
	req, err := http.NewRequest("GET", u, nil)
	req.SetBasicAuth(sb.AuthUsername, sb.AuthPassword)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to fetch catalog from %s\n%v\n", u, resp)
		fmt.Printf("err: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	var catalog model.Catalog
	err = utils.ResponseBodyToObject(resp, &catalog)
	if err != nil {
		fmt.Printf("Failed to unmarshal catalog: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	sbRes := model.CreateServiceBrokerResponse{
		Metadata: model.ServiceBrokerMetadata{
			GUID:      sb.GUID,
			CreatedAt: time.Unix(sb.Created, 0).Format(time.RFC3339),
			URL:       sb.SelfURL,
		},
		Entity: model.ServiceBrokerEntity{
			Name:         sb.Name,
			BrokerURL:    sb.BrokerURL,
			AuthUsername: sb.AuthUsername,
		},
	}

	err = c.storage.AddBroker(&sb, &catalog)
	utils.WriteResponse(w, 200, sbRes)
}

func (c *Controller) DeleteServiceBroker(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractVarFromRequest(r, "broker_id")
	fmt.Printf("DeleteServiceBroker: %s\n", id)

	err := c.storage.DeleteBroker(id)
	if err != nil {
		fmt.Printf("Got Error: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	w.WriteHeader(204)
}

//
// Service Instances.
//

func (c *Controller) ListServiceInstances(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 400, "IMPLEMENT ME")
}

func (c *Controller) GetServiceInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting Service Instance\n")
	id := utils.ExtractVarFromRequest(r, "service_id")

	si, err := c.storage.GetService(id)
	if err != nil {
		fmt.Printf("Couldn't fetch the service: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	utils.WriteResponse(w, 200, si)
}

func (c *Controller) CreateServiceInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Creating Service Instance\n")

	var req CreateServiceInstanceRequest
	err := utils.BodyToObject(r, &req)
	if err != nil {
		fmt.Printf("Error unmarshaling: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	serviceID, err := c.getServiceID(req.ServicePlanGUID)
	if err != nil {
		fmt.Printf("Error fetching service ID: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// Then actually make the request to reify the service instance
	createReq := &ServiceInstanceRequest{
		ServiceID:  serviceID,
		PlanID:     req.ServicePlanGUID,
		Parameters: req.Parameters,
	}

	jsonBytes, err := json.Marshal(createReq)
	if err != nil {
		fmt.Printf("Failed to marshal: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	instanceID := uuid.NewV4().String()

	broker, err := c.getBroker(serviceID)
	if err != nil {
		fmt.Printf("Error fetching service: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	url := fmt.Sprintf(CREATE_SERVICE_INSTANCE_FMT_STR, broker.BrokerURL, instanceID)

	// TODO: Handle the auth
	createHttpReq, err := http.NewRequest("PUT", url, bytes.NewReader(jsonBytes))
	client := &http.Client{}
	fmt.Printf("Doing a request to: %s\n", url)
	resp, err := client.Do(createHttpReq)
	if err != nil {
		fmt.Printf("Failed to PUT: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	defer resp.Body.Close()

	// TODO: Align this with the actual response model.
	si := model.ServiceInstance{}
	err = utils.ResponseBodyToObject(resp, &si)

	si.ID = instanceID
	si.ServiceID = serviceID
	si.PlanID = req.ServicePlanGUID

	c.storage.AddService(&si)
	utils.WriteResponse(w, 200, si)
}

func (c *Controller) DeleteServiceInstance(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 400, "IMPLEMENT ME")
}

func (c *Controller) ListServiceBindings(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 400, "IMPLEMENT ME")
}

func (c *Controller) GetServiceBinding(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting Service Binding\n")
	id := utils.ExtractVarFromRequest(r, "binding_id")

	b, err := c.storage.GetServiceBinding(id)
	if err != nil {
		fmt.Printf("%#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	utils.WriteResponse(w, 400, b)
}

func (c *Controller) CreateServiceBinding(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Creating Service Binding\n")

	var req CreateServiceBindingRequest
	err := utils.BodyToObject(r, &req)
	if err != nil {
		fmt.Printf("Error unmarshaling: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	si, err := c.storage.GetService(req.ServiceInstanceGUID)
	if err != nil {
		fmt.Printf("Error fetching service ID %s: %v\n", req.ServiceInstanceGUID, err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// Then actually make the request to create the binding
	createReq := &BindingRequest{
		ServiceID:  si.ServiceID,
		PlanID:     si.PlanID,
		Parameters: req.Parameters,
	}

	jsonBytes, err := json.Marshal(createReq)
	if err != nil {
		fmt.Printf("Failed to marshal: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	bindingID := uuid.NewV4().String()

	broker, err := c.getBroker(si.ServiceID)
	if err != nil {
		fmt.Printf("Error fetching service: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	url := fmt.Sprintf(BIND_FMT_STR, broker.BrokerURL, si.ID, bindingID)

	// TODO: Handle the auth
	createHttpReq, err := http.NewRequest("PUT", url, bytes.NewReader(jsonBytes))
	client := &http.Client{}
	fmt.Printf("Doing a request to: %s\n", url)
	resp, err := client.Do(createHttpReq)
	if err != nil {
		fmt.Printf("Failed to PUT: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	defer resp.Body.Close()

	sbr := model.CreateServiceBindingResponse{}
	err = utils.ResponseBodyToObject(resp, &sbr)
	if err != nil {
		fmt.Printf("Failed to unmarshal: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// TODO: get broker to actually return these values as part of response.
	sb := model.ServiceBinding{
		ID:                bindingID,
		ServiceInstanceID: si.ID,
		ServiceID:         si.ServiceID,
		ServicePlanID:     si.PlanID,
	}

	c.storage.AddServiceBinding(&sb, &sbr.Credentials)
	utils.WriteResponse(w, 200, sb)
}

func (c *Controller) DeleteServiceBinding(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 400, "IMPLEMENT ME")
}

func (c *Controller) getServiceID(planID string) (string, error) {
	i, err := c.storage.GetInventory()
	if err != nil {
		return "", err
	}

	for _, s := range i.Services {
		for _, p := range s.Plans {
			if strings.Compare(planID, p.ID) == 0 {
				return s.ID, nil
			}
		}
	}

	return "", fmt.Errorf("Plan ID %s was not found", planID)
}

func (c *Controller) getBroker(serviceID string) (*model.ServiceBroker, error) {
	broker, err := c.storage.GetBrokerByService(serviceID)
	if err != nil {
		return nil, err
	}

	return broker, nil
}

// This is what we get sent to us
type CreateServiceInstanceRequest struct {
	Name            string                 `json:"name"`
	ServicePlanGUID string                 `json:"service_plan_guid"`
	Parameters      map[string]interface{} `json:"parameters,omitempty"`
}

type CreateServiceBindingRequest struct {
	ServiceInstanceGUID string                 `json:"service_instance_guid"`
	Parameters          map[string]interface{} `json:"parameters,omitempty"`
}

type ServiceInstanceRequest struct {
	OrgID             string                 `json:"organization_guid,omitempty"`
	PlanID            string                 `json:"plan_id,omitempty"`
	ServiceID         string                 `json:"service_id,omitempty"`
	SpaceID           string                 `json:"space_id,omitempty"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
	AcceptsIncomplete bool                   `json:"accepts_incomplete,omitempty"`
}

type BindingRequest struct {
	AppGUID      string                 `json:"app_guid,omitempty"`
	PlanID       string                 `json:"plan_id,omitempty"`
	ServiceID    string                 `json:"service_id,omitempty"`
	BindResource map[string]interface{} `json:"bind_resource,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}
