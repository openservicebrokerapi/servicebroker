package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	model "github.com/servicebroker/servicebroker/k8s/service_controller/model"
	"github.com/servicebroker/servicebroker/k8s/service_controller/utils"
	sbmodel "github.com/servicebroker/servicebroker/model/service_broker"
	scmodel "github.com/servicebroker/servicebroker/model/service_controller"
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

func (c *Controller) Services(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Services\n")

	i, err := c.storage.GetServices()
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
	var sbReq scmodel.CreateServiceBrokerRequest
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

	// TODO: the model from SB is fetched and stored directly as the one in the SC model (which the
	// storage operates on). We should convert it from the SB model to SC model before storing.
	var catalog model.Catalog
	err = utils.ResponseBodyToObject(resp, &catalog)
	if err != nil {
		fmt.Printf("Failed to unmarshal catalog: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	sbRes := scmodel.CreateServiceBrokerResponse{
		Metadata: scmodel.ServiceBrokerMetadata{
			GUID:      sb.GUID,
			CreatedAt: time.Unix(sb.Created, 0).Format(time.RFC3339),
			URL:       sb.SelfURL,
		},
		Entity: scmodel.ServiceBrokerEntity{
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
	si, err := c.storage.ListServices()
	if err != nil {
		fmt.Printf("Couldn't list services: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	var instances []*model.ServiceInstance
	for _, i := range si {
		instances = append(instances, i)
	}

	utils.WriteResponse(w, 200, instances)
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

	var req scmodel.CreateServiceInstanceRequest
	if err := utils.BodyToObject(r, &req); err != nil {
		fmt.Printf("Error unmarshaling CreateServiceInstanceRequest: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	serviceID, err := c.getServiceID(req.ServicePlanGUID)
	if err != nil {
		err = fmt.Errorf("Error fetching service ID: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	si, err := c.getServiceInstanceByName(req.Name)
	if err != nil {
		err = fmt.Errorf("Error fetching service ID: %v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	existed := (si != nil)
	if si == nil {
		si = &model.ServiceInstance{
			ID: uuid.NewV4().String(),
		}
	}

	si.Name = req.Name
	si.PlanGUID = req.ServicePlanGUID
	si.SpaceGUID = req.SpaceID
	si.Parameters = req.Parameters
	si.Tags = req.Tags
	si.AcceptsIncomplete = r.URL.Query().Get("accepts_incomplete") == "true"

	/* WHAT IS THIS? -Dug

	// Binding data is passed to the service broker right now as part of the
	// parameters in the form:
	//
	// parameters:
	//   bindings:
	//     <service-name>:
	//       <credential>
	if si.Bindings != nil {
		if req.Parameters == nil {
			req.Parameters = make(map[string]interface{})
		}
		req.Parameters["bindings"] = siData.Bindings
	}
	*/

	// Then actually make the request to reify the service instance
	createReq := &sbmodel.CreateServiceInstanceRequest{
		OrgID:             "",
		PlanID:            req.ServicePlanGUID,
		ServiceID:         serviceID,
		Parameters:        si.Parameters,
		AcceptsIncomplete: si.AcceptsIncomplete,
	}

	jsonBytes, err := json.Marshal(createReq)
	if err != nil {
		err = fmt.Errorf("Failed to marshal CreateRequest: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	broker, err := c.getBroker(serviceID)
	if err != nil {
		fmt.Printf("Error fetching service: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	url := fmt.Sprintf(CREATE_SERVICE_INSTANCE_FMT_STR, broker.BrokerURL, si.ID)

	// TODO: Handle the auth
	createHttpReq, err := http.NewRequest("PUT", url, bytes.NewReader(jsonBytes))
	client := &http.Client{}
	fmt.Printf("Doing a request to: %s\n", url)
	resp, err := client.Do(createHttpReq)
	defer resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("Failed to PUT: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	sir := sbmodel.CreateServiceInstanceResponse{}
	if err = utils.ResponseBodyToObject(resp, &sir); err != nil {
		err = fmt.Errorf("Failed to PUT: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	fmt.Printf("Response: %#v\n", sir)

	si.DashboardURL = sir.DashboardURL
	if sir.LastOperation != nil {
		si.LastOperation = &model.LastOperation{
			State:       sir.LastOperation.State,
			Description: sir.LastOperation.Description,
		}
	}

	if existed {
		c.storage.SetService(si)
	} else {
		c.storage.AddService(si)
	}

	utils.WriteResponse(w, 200, si)
}

func (c *Controller) DeleteServiceInstance(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 400, "IMPLEMENT ME")
}

func (c *Controller) ListServiceBindings(w http.ResponseWriter, r *http.Request) {
	l, err := c.storage.ListServiceBindings()
	if err != nil {
		fmt.Printf("Got Error: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	utils.WriteResponse(w, 200, l)
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

	utils.WriteResponse(w, 200, b)
}

func (c *Controller) CreateServiceBinding(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Creating Service Binding\n")

	var req scmodel.CreateServiceBindingRequest
	err := utils.BodyToObject(r, &req)
	if err != nil {
		fmt.Printf("Error unmarshaling: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// Validate that from service has not been instantiated yet.
	fromSI, err := c.getServiceInstanceByName(req.FromServiceInstanceName)
	if err != nil {
		fromSI = &model.ServiceInstance{
			Name:     req.FromServiceInstanceName,
			ID:       uuid.NewV4().String(),
			Bindings: make(map[string]*interface{}), // Credentials
		}
		c.storage.AddService(fromSI)
	}

	if fromSI.Service.ID != "" {
		err = fmt.Errorf("Cannot bind from instantiated service: %s (%s)", req.FromServiceInstanceName, fromSI.ID)
		utils.WriteResponse(w, 400, err)
		return
	}

	// Get instance information for service being bound to.
	si, err := c.storage.GetService(req.ServiceInstanceGUID)
	if err != nil {
		fmt.Printf("Error fetching service ID %s: %v\n", req.ServiceInstanceGUID, err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// Then actually make the request to create the binding
	createReq := &sbmodel.CreateServiceBindingRequest{
		ServiceID:  si.Service.ID,
		PlanID:     si.PlanGUID,
		Parameters: req.Parameters,
	}

	jsonBytes, err := json.Marshal(createReq)
	if err != nil {
		fmt.Printf("Failed to marshal: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	bindingID := uuid.NewV4().String()

	broker, err := c.getBroker(si.Service.ID)
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
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("Failed to PUT: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	sbr := scmodel.CreateServiceBindingResponse{}
	err = utils.ResponseBodyToObject(resp, &sbr)
	if err != nil {
		fmt.Printf("Failed to unmarshal: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// TODO: get broker to actually return these values as part of response.
	sb := model.ServiceBinding{
		ID: bindingID,
		FromServiceInstanceName: req.FromServiceInstanceName,
		ServiceInstanceGUID:     req.ServiceInstanceGUID,
		Parameters:              req.Parameters,
	}

	c.storage.AddServiceBinding(&sb, &sbr.Credentials)

	// Set binding credential information in from service instance.
	serviceName, err := c.getServiceName(si.ID)
	if err != nil {
		fmt.Printf("Error retrieving service name: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	fromSI.Bindings[serviceName] = &sbr.Credentials
	c.storage.SetService(fromSI)

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

func (c *Controller) getServiceName(instanceId string) (string, error) {
	si, err := c.storage.GetService(instanceId)
	if err != nil {
		return "", err
	}

	i, err := c.storage.GetInventory()
	if err != nil {
		return "", err
	}

	for _, s := range i.Services {
		if strings.Compare(si.Service.ID, s.ID) == 0 {
			return s.Name, nil
		}
	}

	return "", fmt.Errorf("Service ID %s was not found for instance %s", si.Service.ID, instanceId)
}

func (c *Controller) getBroker(serviceID string) (*model.ServiceBroker, error) {
	broker, err := c.storage.GetBrokerByService(serviceID)
	if err != nil {
		return nil, err
	}

	return broker, nil
}

func (c *Controller) getServiceInstanceByName(name string) (*model.ServiceInstance, error) {
	siList, err := c.storage.ListServices()
	if err != nil {
		return nil, err
	}
	for _, si := range siList {
		if strings.Compare(si.Name, name) == 0 {
			return si, nil
		}
	}

	return nil, nil
}
