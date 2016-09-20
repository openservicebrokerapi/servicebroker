package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	storage model.ServiceStorage
}

func CreateController(storage model.ServiceStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}

//
// Inventory.
//

func (c *Controller) Inventory(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Inventory\n")

	serviceIDs, err := c.storage.ListServices()
	if err != nil {
		err = fmt.Errorf("Error getting services: %#v", err)
		utils.WriteResponse(w, 500, err)
		return
	}
	gcr := scmodel.GetCatalogResponse{}

	for _, svcID := range serviceIDs {
		svc, err := c.storage.GetService(svcID)
		if err != nil {
			err = fmt.Errorf("Error getting service %q: %#v", svcID, err)
			utils.WriteResponse(w, 500, err)
			return
		}
		if svc == nil {
			utils.WriteResponse(w, 500, fmt.Errorf("Can't find service %q", svcID))
			return
		}

		newSvc := scmodel.Service{
			Name:            svc.Name,
			ID:              svc.ID,
			Description:     svc.Description,
			Bindable:        svc.Bindable,
			PlanUpdateable:  svc.PlanUpdateable,
			Tags:            svc.Tags,
			Requires:        svc.Requires,
			Metadata:        svc.Metadata,
			Plans:           []scmodel.ServicePlan{},
			DashboardClient: nil,
		}

		for _, planID := range svc.Plans {
			plan, err := c.storage.GetPlan(planID)
			if err != nil {
				err = fmt.Errorf("Error getting plan %q: %#v", planID, err)
				utils.WriteResponse(w, 500, err)
				return
			}
			if plan == nil {
				utils.WriteResponse(w, 500, fmt.Errorf("Can't find plan %q", planID))
				return
			}

			newPlan := scmodel.ServicePlan{
				Name:        plan.Name,
				ID:          plan.ID,
				Description: plan.Description,
				Metadata:    plan.Metadata,
				Free:        plan.Free,
				// Schemas:     plan.Schemas,
			}
			newSvc.Plans = append(newSvc.Plans, newPlan)
		}
		gcr.Services = append(gcr.Services, newSvc)
	}

	utils.WriteResponse(w, 200, gcr)
}

func (c *Controller) Services(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Services\n")

	serviceIDs, err := c.storage.ListServices()
	if err != nil {
		err = fmt.Errorf("Error getting services: %#v", err)
		utils.WriteResponse(w, 500, err)
		return
	}

	services := []*scmodel.Service{}
	for _, serviceID := range serviceIDs {
		svc, err := c.storage.GetService(serviceID)
		if err != nil {
			err = fmt.Errorf("Error getting service %q: %#v", serviceID, err)
			utils.WriteResponse(w, 500, err)
			return
		}
		if svc == nil {
			utils.WriteResponse(w, 500, fmt.Errorf("Can't find service %q", serviceID))
			return
		}

		service := scmodel.Service{
			Name:           svc.Name,
			ID:             svc.ID,
			Description:    svc.Description,
			Bindable:       svc.Bindable,
			PlanUpdateable: svc.PlanUpdateable,
			Tags:           svc.Tags,
			Requires:       svc.Requires,

			Metadata:        svc.Metadata,
			Plans:           []scmodel.ServicePlan{},
			DashboardClient: svc.DashboardClient,
		}

		for _, planID := range svc.Plans {
			plan, err := c.storage.GetPlan(planID)
			if err != nil {
				err = fmt.Errorf("Error getting plan %q: %#v", planID, err)
				utils.WriteResponse(w, 500, err)
				return
			}
			if plan == nil {
				utils.WriteResponse(w, 500, fmt.Errorf("Can't find plan %q", planID))
				return
			}

			service.Plans = append(service.Plans, scmodel.ServicePlan{
				Name:        plan.Name,
				ID:          plan.ID,
				Description: plan.Description,
				Metadata:    plan.Metadata,
				Free:        plan.Free,
				Schemas: scmodel.Schemas{
					Instance: scmodel.Schema{
						Inputs:  plan.Schemas.Instance.Inputs,
						Outputs: plan.Schemas.Instance.Outputs,
					},
					Binding: scmodel.Schema{
						Inputs:  plan.Schemas.Instance.Inputs,
						Outputs: plan.Schemas.Instance.Outputs,
					},
				},
			})
		}

		services = append(services, &service)
	}

	utils.WriteResponse(w, 200, services)
}

//
// Service Broker.
//

func (c *Controller) ListServiceBrokers(w http.ResponseWriter, r *http.Request) {
	brokerIDs, err := c.storage.ListBrokers()
	if err != nil {
		err = fmt.Errorf("Error getting brokers: %#v", err)
		utils.WriteResponse(w, 500, err)
		return
	}

	brokers := []scmodel.ServiceBroker{}
	for _, bID := range brokerIDs {
		broker, err := c.storage.GetBroker(bID)
		if err != nil {
			err = fmt.Errorf("Error getting bindings: %#v", err)
			utils.WriteResponse(w, 500, err)
			return
		}
		if broker == nil {
			utils.WriteResponse(w, 500, fmt.Errorf("Can't find broker %q", bID))
			return
		}

		brokers = append(brokers, scmodel.ServiceBroker{
			GUID:         broker.ID,
			Name:         broker.Name,
			BrokerURL:    broker.BrokerURL,
			AuthUsername: broker.AuthUsername,
			AuthPassword: broker.AuthPassword,

			Created: broker.Created,
			Updated: broker.Updated,
			SelfURL: broker.SelfURL,
		})
	}

	utils.WriteResponse(w, 200, brokers)
}

func (c *Controller) GetServiceBroker(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractVarFromRequest(r, "broker_id")
	fmt.Printf("GetServiceBroker: %s\n", id)

	b, err := c.storage.GetBroker(id)
	if err != nil {
		err = fmt.Errorf("Error finding broker %q: %#v", id, err)
		utils.WriteResponse(w, 500, err)
		return
	}
	if b == nil {
		utils.WriteResponse(w, 404, fmt.Errorf("Can't find service broker: %s", id))
		return
	}

	utils.WriteResponse(w, 200, b)
}

func (c *Controller) CreateServiceBroker(w http.ResponseWriter, r *http.Request) {
	var sbReq scmodel.CreateServiceBrokerRequest
	err := utils.BodyToObject(r, &sbReq)
	if err != nil {
		err = fmt.Errorf("Error unmarshaling CreateSvcBrokerResp: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	newGUID := uuid.NewV4().String()

	// Fetch the catalog from the broker
	u := fmt.Sprintf(CATALOG_URL_FMT_STR, sbReq.BrokerURL)
	req, _ := http.NewRequest("GET", u, nil)
	req.SetBasicAuth(sbReq.AuthUsername, sbReq.AuthPassword)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("Error Getting catalog from %s\n%v\n:%#v", u, resp, err)
		utils.WriteResponse(w, 400, err)
		return
	}

	var catalog sbmodel.GetCatalogResponse
	if err = utils.ResponseBodyToObject(resp, &catalog); err != nil {
		err = fmt.Errorf("Failed to unmarshal catalog: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	sb := model.ServiceBroker{
		ID:           newGUID,
		Name:         sbReq.Name,
		BrokerURL:    sbReq.BrokerURL,
		AuthUsername: sbReq.AuthUsername,
		AuthPassword: sbReq.AuthPassword,

		Created:  time.Now().Unix(),
		Updated:  0,
		SelfURL:  "/v2/service_brokers/" + newGUID,
		Services: []string{},
	}

	for _, catSvc := range catalog.Services {
		newSvc := model.Service{
			Name:           catSvc.Name,
			ID:             catSvc.ID,
			Description:    catSvc.Description,
			Bindable:       catSvc.Bindable,
			PlanUpdateable: catSvc.PlanUpdateable,
			Tags:           catSvc.Tags,
			Requires:       catSvc.Requires,

			Metadata:        catSvc.Metadata,
			DashboardClient: catSvc.DashboardClient,

			ServiceBroker: sb.ID,
			Plans:         nil,
			Instances:     nil,
		}

		for _, catPlan := range catSvc.Plans {
			newPlan := model.ServicePlan{
				Name:        catPlan.Name,
				ID:          catPlan.ID,
				Description: catPlan.Description,
				Metadata:    catPlan.Metadata,
				Free:        catPlan.Free,
				// Schemas:     catPlan.Schemas,  // need to convert it

				Service: newSvc.ID,
			}
			err = c.storage.AddPlan(&newPlan)
			if err != nil {
				err = fmt.Errorf("Error adding plan %s: %#v", newPlan.ID, err)
				utils.WriteResponse(w, 500, err)
				return
			}

			newSvc.Plans = append(newSvc.Plans, newPlan.ID)
		}
		err = c.storage.AddService(&newSvc)
		if err != nil {
			err = fmt.Errorf("Error adding service %s: %#v", newSvc.ID, err)
			utils.WriteResponse(w, 500, err)
			return
		}
		sb.Services = append(sb.Services, newSvc.ID)
	}

	if err = c.storage.AddBroker(&sb); err != nil {
		err = fmt.Errorf("Error saving new broker: %#v", err)
		utils.WriteResponse(w, 500, err)
		return
	}

	sbRes := scmodel.CreateServiceBrokerResponse{
		Metadata: scmodel.ServiceBrokerMetadata{
			GUID:      sb.ID,
			CreatedAt: time.Unix(sb.Created, 0).Format(time.RFC3339),
			URL:       sb.SelfURL,
		},
		Entity: scmodel.ServiceBrokerEntity{
			Name:         sb.Name,
			BrokerURL:    sb.BrokerURL,
			AuthUsername: sb.AuthUsername,
		},
	}

	utils.WriteResponse(w, 200, sbRes)
}

func (c *Controller) DeleteServiceBroker(w http.ResponseWriter, r *http.Request) {
	id := utils.ExtractVarFromRequest(r, "broker_id")
	fmt.Printf("DeleteServiceBroker: %s\n", id)

	broker, err := c.storage.GetBroker(id)
	if err != nil {
		err = fmt.Errorf("Error getting broker %q: %#v", id, err)
		utils.WriteResponse(w, 500, err)
		return
	}
	if broker == nil {
		err = fmt.Errorf("Can't find broker %q", id)
		utils.WriteResponse(w, 404, err)
		return
	}

	// TODO consider moving this out of this http layer

	tmpServices := []string{} // Play it safe
	copy(tmpServices, broker.Services)

	for _, serviceID := range tmpServices {
		service, err := c.storage.GetService(serviceID)
		if err != nil {
			err = fmt.Errorf("Error getting service %s: %#v", serviceID, err)
			utils.WriteResponse(w, 500, err)
			return
		}
		if service == nil {
			err = fmt.Errorf("Can't find service %s", serviceID)
			utils.WriteResponse(w, 500, err)
			return
		}

		c.storage.DeleteService(service.ID)
	}
	broker.Services = []string{}
	c.storage.DeleteBroker(broker.ID)

	w.WriteHeader(204)
}

//
// Service Instances.
//

func (c *Controller) ListServiceInstances(w http.ResponseWriter, r *http.Request) {
	instanceIDs, err := c.storage.ListInstances()
	if err != nil {
		err = fmt.Errorf("Couldn't list instances: %#v\n", err)
		utils.WriteResponse(w, 500, err)
		return
	}

	instances := []scmodel.ServiceInstance{}
	for _, iID := range instanceIDs {
		instance, err := c.storage.GetInstance(iID)
		if err != nil {
			err = fmt.Errorf("Error getting instance: %#v", err)
			utils.WriteResponse(w, 500, err)
			return
		}
		if instance == nil {
			utils.WriteResponse(w, 500, fmt.Errorf("Can't find instance %q", iID))
			return
		}

		instances = append(instances, scmodel.ServiceInstance{
			Name: instance.Name,
			// Credentials: ??? TODO FIX
			ServicePlanGUID: instance.Plan,
			SpaceGUID:       instance.SpaceID,
			DashboardURL:    instance.DashboardURL,
			Type:            "managed_service_instance",
			LastOperation: scmodel.LastOperation{
				Type:        "create",
				State:       instance.LastOperation.State,
				Description: instance.LastOperation.Description,
				UpdatedAt:   instance.LastOperation.UpdatedAt,
			},
			SpaceURL:       "",
			ServicePlanURL: "",
			RoutesURL:      "",
			Tags:           instance.Tags,

			Parameters: instance.Parameters,

			ID:        instance.ID,
			ServiceID: instance.Service,
		})
	}

	utils.WriteResponse(w, 200, instances)
}

func (c *Controller) GetServiceInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting Service Instance\n")
	id := utils.ExtractVarFromRequest(r, "service_id")

	si, err := c.storage.GetInstance(id)
	if err != nil {
		err = fmt.Errorf("Error finding instance %q: %#v", id, err)
		utils.WriteResponse(w, 500, err)
		return
	}
	if si == nil {
		utils.WriteResponse(w, 404, fmt.Errorf("Can't find service instance: %s", id))
		return
	}

	res := scmodel.ServiceInstance{
		Name:            si.Name,
		Credentials:     "", // TODO - old?
		ServicePlanGUID: "",
		SpaceGUID:       si.SpaceID,
		DashboardURL:    si.DashboardURL,
		Type:            "managed_service_instance",
		// LastOperation:   // TODO - get real value
		SpaceURL:       "",
		ServicePlanURL: "",
		RoutesURL:      "",
		Tags:           si.Tags,

		Parameters: si.Parameters,

		ID:        si.ID,
		ServiceID: si.Service,
	}

	utils.WriteResponse(w, 200, &res)
}

func (c *Controller) CreateServiceInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Creating Service Instance\n")

	var req scmodel.CreateServiceInstanceRequest
	if err := utils.BodyToObject(r, &req); err != nil {
		fmt.Printf("Error unmarshaling CreateServiceInstanceRequest: %v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	plan, err := c.storage.GetPlan(req.ServicePlanGUID)
	if err != nil {
		err = fmt.Errorf("Error getting plan %q: %#v\n", req.ServicePlanGUID, err)
		utils.WriteResponse(w, 500, err)
		return
	}
	if plan == nil {
		err = fmt.Errorf("Can't find plan %q", req.ServicePlanGUID)
		utils.WriteResponse(w, 404, err)
		return
	}

	service, err := c.storage.GetService(plan.Service)
	if err != nil {
		err = fmt.Errorf("Error fetching service %q: %#v", plan.Service, err)
		utils.WriteResponse(w, 500, err)
		return
	}
	if service == nil {
		utils.WriteResponse(w, 500, fmt.Errorf("Can't find service %q", plan.Service))
		return
	}

	// TODO: get rid of getServiceInstanceByName - maybe???
	si, err := c.getServiceInstanceByName(req.Name)
	if err != nil {
		err = fmt.Errorf("Error fetching service ID: %v", err)
		utils.WriteResponse(w, 500, err)
		return
	}

	if si == nil {
		si = &model.ServiceInstance{
			ID: uuid.NewV4().String(),
		}
	}

	si.Name = req.Name
	si.Plan = req.ServicePlanGUID
	si.SpaceID = req.SpaceID
	si.Parameters = req.Parameters
	si.Tags = req.Tags
	si.AcceptsIncomplete = r.URL.Query().Get("accepts_incomplete") == "true"
	// si.DashboardURL = ...
	// si.LastOperion = ...
	si.Service = service.ID
	// si.Bindings = ...

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
		ServiceID:         plan.Service,
		Parameters:        req.Parameters,
		AcceptsIncomplete: si.AcceptsIncomplete,
	}

	jsonBytes, err := json.Marshal(createReq)
	if err != nil {
		err = fmt.Errorf("Failed to marshal CreateRequest: %#v", err)
		utils.WriteResponse(w, 500, err)
		return
	}

	broker, err := c.storage.GetBroker(service.ServiceBroker)
	if err != nil {
		err = fmt.Errorf("Error fetching broker %q: %#v", service.ServiceBroker, err)
		utils.WriteResponse(w, 500, err)
		return
	}
	if broker == nil {
		utils.WriteResponse(w, 404, fmt.Errorf("Can't find service broker for service: %s", service.ServiceBroker))
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
		err = fmt.Errorf("Failed to PUT to Service Broker: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	sir := sbmodel.CreateServiceInstanceResponse{}
	if err = utils.ResponseBodyToObject(resp, &sir); err != nil {
		err = fmt.Errorf("Failed to PUT to Service Broker: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}
	fmt.Printf("Response: %#v\n", sir)

	si.DashboardURL = sir.DashboardURL
	if sir.LastOperation != nil {
		si.LastOperation = model.LastOperation{
			UpdatedAt:                "", // TODO - complete
			State:                    sir.LastOperation.State,
			Description:              sir.LastOperation.Description,
			AsyncPollIntervalSeconds: 0,
		}
	}

	si.Service = service.ID

	if service.Instances == nil {
		service.Instances = map[string]string{}
	}
	service.Instances[si.ID] = si.ID // TODO: do we really need map??

	err = c.storage.AddInstance(si)
	if err != nil {
		err = fmt.Errorf("Error adding instance %s: %#v", si.ID, err)
		utils.WriteResponse(w, 500, err)
		return
	}

	err = c.storage.SetService(service)
	if err != nil {
		err = fmt.Errorf("Error saving service %s: %#v", service.ID, err)
		utils.WriteResponse(w, 500, err)
		return
	}

	res := scmodel.ServiceInstance{
		Name: si.Name,
		// Credentials: ??? TODO FIX
		ServicePlanGUID: si.Plan,
		SpaceGUID:       si.SpaceID,
		DashboardURL:    si.DashboardURL,
		Type:            "managed_service_instance",
		LastOperation: scmodel.LastOperation{
			Type:        "create",
			State:       si.LastOperation.State,
			Description: si.LastOperation.Description,
			UpdatedAt:   si.LastOperation.UpdatedAt,
		},
		SpaceURL:       "",
		ServicePlanURL: "",
		RoutesURL:      "",
		Tags:           si.Tags,

		Parameters: si.Parameters,

		ID:        si.ID,
		ServiceID: si.Service,
	}

	utils.WriteResponse(w, 200, res)
}

func (c *Controller) DeleteServiceInstance(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 500, "IMPLEMENT ME")
}

func (c *Controller) ListServiceBindings(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Listing Service Binding\n")
	bindingIDs, err := c.storage.ListBindings()
	if err != nil {
		err = fmt.Errorf("Error getting bindings: %#v", err)
		utils.WriteResponse(w, 500, err)
		return
	}

	daBindings := make([]scmodel.ServiceBinding, 0, len(bindingIDs))
	for _, bID := range bindingIDs {
		b, err := c.storage.GetBinding(bID)
		if err != nil {
			err = fmt.Errorf("Error getting binding %q: %#v", bID, err)
			utils.WriteResponse(w, 500, err)
			return
		}
		if b == nil {
			utils.WriteResponse(w, 500, fmt.Errorf("Can't find binding %q", bID))
			return
		}

		binding := scmodel.ServiceBinding{
			ID:                  b.ID,
			AppName:             b.AppName,
			ServiceInstanceName: b.ServiceInstanceName,
			ServiceInstanceGUID: b.ServiceInstanceID,
			Parameters:          b.Parameters,
			Credentials:         b.Credentials,
		}
		daBindings = append(daBindings, binding)
	}

	utils.WriteResponse(w, 200, daBindings)
}

func (c *Controller) GetServiceBinding(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Getting Service Binding\n")

	id := utils.ExtractVarFromRequest(r, "binding_id")
	b, err := c.storage.GetBinding(id)
	if err != nil {
		err = fmt.Errorf("Error finding service binding %q: %#v", id, err)
		utils.WriteResponse(w, 500, err)
		return
	}
	if b == nil {
		utils.WriteResponse(w, 404, fmt.Errorf("Can't find binding: %s", id))
	}

	utils.WriteResponse(w, 200, b)
}

func (c *Controller) CreateServiceBinding(w http.ResponseWriter, r *http.Request) {
	var req scmodel.CreateServiceBindingRequest
	err := utils.BodyToObject(r, &req)
	if err != nil {
		err = fmt.Errorf("Error unmarshaling: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	fmt.Printf("REQ:\n%#v\n\n", req)

	// Validate that from service has not been instantiated yet.
	instance, err := c.getServiceInstanceByName(req.ServiceInstanceName)
	if err != nil {
		err = fmt.Errorf("Error finding service %q: %#v", req.ServiceInstanceName, err)
		utils.WriteResponse(w, 500, err)
		return
	}
	if instance == nil {
		err = fmt.Errorf("Can't find service %q: %#v", req.ServiceInstanceName, err)
		utils.WriteResponse(w, 404, err)
		return
	}

	// Then actually make the request to create the binding
	createReq := &sbmodel.CreateServiceBindingRequest{
		AppGUID:      req.AppName,
		PlanID:       instance.Plan,
		ServiceID:    instance.Service,
		BindResource: nil,
		Parameters:   req.Parameters,
	}

	jsonBytes, err := json.Marshal(createReq)
	if err != nil {
		err = fmt.Errorf("Failed to marshal request: %#v\n", err)
		utils.WriteResponse(w, 500, err)
		return
	}

	service, err := c.storage.GetService(instance.Service)
	if err != nil {
		err = fmt.Errorf("Error fetching service %q: %#v", instance.Service, err)
		utils.WriteResponse(w, 500, err)
		return
	}
	if service == nil {
		utils.WriteResponse(w, 404, fmt.Errorf("Can't find service: %s", instance.Service))
		return
	}

	broker, err := c.storage.GetBroker(service.ServiceBroker)
	if err != nil {
		err = fmt.Errorf("Error fetching broker %q: %#v", service.ServiceBroker, err)
		utils.WriteResponse(w, 500, err)
		return
	}
	if broker == nil {
		utils.WriteResponse(w, 404, fmt.Errorf("Can't find service broker: %s", instance.Service))
		return
	}

	bindingID := uuid.NewV4().String()

	url := fmt.Sprintf(BIND_FMT_STR, broker.BrokerURL, instance.ID, bindingID)

	// TODO: Handle the auth
	createHttpReq, err := http.NewRequest("PUT", url, bytes.NewReader(jsonBytes))
	client := &http.Client{}
	fmt.Printf("Doing a request to: %s\n", url)
	resp, err := client.Do(createHttpReq)
	defer resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("Failed to PUT to Service Broker: %#v", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	sbr := sbmodel.CreateServiceBindingResponse{}
	if err = utils.ResponseBodyToObject(resp, &sbr); err != nil {
		err = fmt.Errorf("Failed to unmarshal response: %#v\n", err)
		utils.WriteResponse(w, 400, err)
		return
	}

	// TODO: get broker to actually return these values as part of response.
	binding := model.ServiceBinding{
		Credentials:     sbr.Credentials,
		SyslogDrainURL:  sbr.SyslogDrainURL,
		RouteServiceURL: sbr.RouteServiceURL,
		VolumeMounts:    sbr.VolumeMounts,

		AppName:             req.AppName,
		ServiceInstanceName: req.ServiceInstanceName, // Why both name and guid?
		ServiceInstanceID:   req.ServiceInstanceGUID,
		Parameters:          req.Parameters,

		ID: bindingID,
	}
	err = c.storage.AddBinding(&binding)
	if err != nil {
		err = fmt.Errorf("Error adding binding %s: %#v", binding.ID, err)
		utils.WriteResponse(w, 500, err)
		return
	}

	if instance.Bindings == nil {
		instance.Bindings = map[string]string{}
	}
	instance.Bindings[binding.ID] = binding.ID
	err = c.storage.SetInstance(instance)
	if err != nil {
		err = fmt.Errorf("Error saving instance %s: %#v", instance.ID, err)
		utils.WriteResponse(w, 500, err)
		return
	}

	res := scmodel.ServiceBinding{
		ID:                  binding.ID,
		AppName:             binding.AppName,
		ServiceInstanceName: binding.ServiceInstanceName,
		ServiceInstanceGUID: binding.ServiceInstanceID,
		Parameters:          binding.Parameters,
		Credentials:         binding.Credentials,
	}

	utils.WriteResponse(w, 200, res)
}

func (c *Controller) DeleteServiceBinding(w http.ResponseWriter, r *http.Request) {
	utils.WriteResponse(w, 500, "IMPLEMENT ME")
}

func (c *Controller) getServiceInstanceByName(name string) (*model.ServiceInstance, error) {
	instanceIDs, err := c.storage.ListInstances()
	if err != nil {
		return nil, fmt.Errorf("Couldn't list instances: %#v\n", err)
	}

	for _, instanceID := range instanceIDs {
		instance, err := c.storage.GetInstance(instanceID)
		if err != nil {
			return nil, fmt.Errorf("Error getting instance %q: %#v", instanceID, err)
		}
		if instance == nil {
			return nil, fmt.Errorf("Can't find instance %q", instanceID)
		}

		if instance.Name == name {
			return instance, nil
		}
	}

	return nil, nil
}
