package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/satori/go.uuid"
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

	sb := scmodel.ServiceBroker{
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
	var catalog scmodel.Catalog
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

	var instances []*scmodel.ServiceInstance
	for _, i := range si {
		instances = append(instances, i.Instance)
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

	utils.WriteResponse(w, 200, si.Instance)
}

func (c *Controller) CreateServiceInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Creating Service Instance\n")

	var req scmodel.CreateServiceInstanceRequest
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

	siData, err := c.getServiceInstanceByName(req.Name)
	if err != nil {
		siData = &scmodel.ServiceInstanceData{Instance: &scmodel.ServiceInstance{
			ID: uuid.NewV4().String(),
		}}
	}
	existed := (err == nil)

	// Binding data is passed to the service broker right now as part of the
	// parameters in the form:
	//
	// parameters:
	//   bindings:
	//     <service-name>:
	//       <credential>
	if siData.Bindings != nil {
		if req.Parameters == nil {
			req.Parameters = make(map[string]interface{})
		}
		req.Parameters["bindings"] = siData.Bindings
	}

	// Then actually make the request to reify the service instance
	createReq := &sbmodel.ServiceInstanceRequest{
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

	broker, err := c.getBroker(serviceID)
	if err != nil {
		fmt.Printf("Error fetching service: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	url := fmt.Sprintf(CREATE_SERVICE_INSTANCE_FMT_STR, broker.BrokerURL, siData.Instance.ID)

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
	si := scmodel.ServiceInstance{}
	err = utils.ResponseBodyToObject(resp, &si)

	si.Name = req.Name
	si.ID = siData.Instance.ID
	si.ServiceID = serviceID
	si.PlanID = req.ServicePlanGUID

	siData.Instance = &si

	if existed {
		c.storage.SetService(siData)
	} else {
		c.storage.AddService(siData)
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
		fromSI = &scmodel.ServiceInstanceData{
			Instance: &scmodel.ServiceInstance{
				Name: req.FromServiceInstanceName,
				ID:   uuid.NewV4().String(),
			},
			Bindings: make(map[string]*scmodel.Credential),
		}
		c.storage.AddService(fromSI)
	}

	if fromSI.Instance.ServiceID != "" {
		err = fmt.Errorf("Cannot bind from instantiated service: %s (%s)", req.FromServiceInstanceName, fromSI.Instance.ID)
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
	createReq := &sbmodel.BindingRequest{
		ServiceID:  si.Instance.ServiceID,
		PlanID:     si.Instance.PlanID,
		Parameters: req.Parameters,
	}

	jsonBytes, err := json.Marshal(createReq)
	if err != nil {
		fmt.Printf("Failed to marshal: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	bindingID := uuid.NewV4().String()

	broker, err := c.getBroker(si.Instance.ServiceID)
	if err != nil {
		fmt.Printf("Error fetching service: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	url := fmt.Sprintf(BIND_FMT_STR, broker.BrokerURL, si.Instance.ID, bindingID)

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

	sbr := scmodel.CreateServiceBindingResponse{}
	err = utils.ResponseBodyToObject(resp, &sbr)
	if err != nil {
		fmt.Printf("Failed to unmarshal: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// TODO: get broker to actually return these values as part of response.
	sb := scmodel.ServiceBinding{
		ID: bindingID,
		FromServiceInstanceName: req.FromServiceInstanceName,
		ServiceInstanceGUID:     req.ServiceInstanceGUID,
		Parameters:              req.Parameters,
	}

	c.storage.AddServiceBinding(&sb, &sbr.Credentials)

	// Set binding credential information in from service instance.
	serviceName, err := c.getServiceName(si.Instance.ID)
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
		if strings.Compare(si.Instance.ServiceID, s.ID) == 0 {
			return s.Name, nil
		}
	}

	return "", fmt.Errorf("Service ID %s was not found for instance %s", si.Instance.ServiceID, instanceId)
}

func (c *Controller) getBroker(serviceID string) (*scmodel.ServiceBroker, error) {
	broker, err := c.storage.GetBrokerByService(serviceID)
	if err != nil {
		return nil, err
	}

	return broker, nil
}

func (c *Controller) getServiceInstanceByName(name string) (*scmodel.ServiceInstanceData, error) {
	siList, err := c.storage.ListServices()
	if err != nil {
		return nil, err
	}
	for _, si := range siList {
		if strings.Compare(si.Instance.Name, name) == 0 {
			return si, nil
		}
	}

	return nil, fmt.Errorf("Service instance %s was not found", name)
}
