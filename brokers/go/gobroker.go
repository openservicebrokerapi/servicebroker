package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	sbmodel "github.com/servicebroker/servicebroker/model/service_broker"
)

type Service struct {
	Name           string      `json:"name"`
	ID             string      `json:"id"`
	Description    string      `json:"description"`
	Bindable       bool        `json:"bindable"`
	PlanUpdateable bool        `json:"plan_updateable"`
	Plans          []*Plan     `json:"plans"`
	Instances      []*Instance `json:"-"`
}

type Plan struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"description"`
	Free        bool   `json:"free"`
}

type Instance struct {
	ID           string
	DashboardURL string
	State        string
	Service      *Service   `json:"-"`
	Bindings     []*Binding `json:"-"`
}

type Binding struct {
	Instance    *Instance
	ID          string
	AppGUID     string
	PlanID      string
	Credentials string
}

var Services = map[string]*Service{}
var Instances = map[string]*Instance{}
var Bindings = map[string]*Binding{}

func init() {
	plan1 := Plan{
		Name:        "freePlan",
		ID:          "ffffffff-0000-0000-0000-000000000000",
		Description: "free is good",
		Free:        true,
	}
	plan2 := Plan{
		Name:        "costlyPlan",
		ID:          "eeeeeeee-1111-1111-1111-111111111111",
		Description: "not so free",
	}

	Service := Service{
		Name:           "myService1",
		ID:             "12345678-abcd-1234-bcde-1234567890ab",
		Description:    "something cool",
		Bindable:       true,
		PlanUpdateable: false,
		Plans:          []*Plan{&plan1, &plan2},
	}

	Services[Service.ID] = &Service
}

func Log(str string) {
	fmt.Printf("%s\n", str)
}

func WriteResponse(w http.ResponseWriter, code int, object interface{}) {
	var data []byte
	var err error

	if str, ok := object.(string); ok {
		data = []byte(str)
	} else if err, ok = object.(error); ok {
		if jerr, ok := err.(*json.SyntaxError); ok {
			data = []byte(fmt.Sprintf("%s - offset: %d", err, jerr.Offset))
		} else {
			data = []byte(err.Error())
		}
	} else {
		data, err = json.Marshal(object)
		if err != nil {
			code = http.StatusInternalServerError
			data = []byte(fmt.Sprintf("%s", err))
		}
	}

	w.WriteHeader(code)
	Log("Response(" + fmt.Sprintf("%d", code) + "): " + string(data) + "\n")
	fmt.Fprintf(w, string(data)+"\n")
}

func getCatalog(w http.ResponseWriter, r *http.Request) {
	// GET /v2/catalog
	Log(r.Method + ": " + r.URL.String())

	catalog := sbmodel.GetCatalogResponse{}

	for _, svc := range Services {
		service := sbmodel.Service{
			Name:           svc.Name,
			ID:             svc.ID,
			Description:    svc.Description,
			Bindable:       svc.Bindable,
			PlanUpdateable: svc.PlanUpdateable,
		}
		for _, daPlan := range svc.Plans {
			plan := sbmodel.ServicePlan{
				Name:        daPlan.Name,
				ID:          daPlan.ID,
				Description: daPlan.Description,
				Free:        daPlan.Free,
			}
			service.Plans = append(service.Plans, plan)
		}
		catalog.Services = append(catalog.Services, &service)
	}

	data, err := json.Marshal(catalog)
	if err != nil {
		WriteResponse(w, 500, err)
		return
	}
	res := string(data) + "\n"
	WriteResponse(w, 200, res)
}

func createInstance(w http.ResponseWriter, r *http.Request) {
	// PUT /v2/service_instances/:instance_id
	Log(r.Method + ": " + r.URL.String())

	var instanceReq sbmodel.CreateServiceInstanceRequest
	err := json.NewDecoder(r.Body).Decode(&instanceReq)
	if err != nil {
		WriteResponse(w, 500, err)
		return
	}

	instanceID := mux.Vars(r)["iid"]
	if instanceID == "" {
		WriteResponse(w, 500, fmt.Sprintf("InstanceID can't be blank"))
		return
	}

	if _, ok := Instances[instanceID]; ok {
		WriteResponse(w, 500, fmt.Sprintf("InstanceID %q is already in use", instanceID))
		return
	}

	service, ok := Services[instanceReq.ServiceID]
	if !ok {
		WriteResponse(w, 500, fmt.Sprintf("Can't find service id: %s", instanceReq.ServiceID))
		return
	}

	instance := Instance{
		ID:           instanceID,
		Service:      service,
		DashboardURL: "http://example.com/dashAwayAll",
		State:        "created",
	}

	Instances[instanceID] = &instance

	instanceRes := sbmodel.CreateServiceInstanceResponse{
		DashboardURL:  instance.DashboardURL,
		LastOperation: nil,
	}

	WriteResponse(w, 201, instanceRes)
}

func createBinding(w http.ResponseWriter, r *http.Request) {
	// PUT /v2/service_instances/:instance_id/service_bindings/:bidning_id
	Log(r.Method + ": " + r.URL.String())

	instanceID := mux.Vars(r)["iid"]
	bindingID := mux.Vars(r)["bid"]

	var bindingReq sbmodel.CreateServiceBindingRequest
	err := json.NewDecoder(r.Body).Decode(&bindingReq)
	if err != nil {
		WriteResponse(w, 500, err)
		return
	}

	service, ok := Services[bindingReq.ServiceID]
	if !ok {
		WriteResponse(w, 500, fmt.Sprintf("Can't find service id: %s", bindingReq.ServiceID))
		return
	}

	instance, ok := Instances[instanceID]
	if !ok {
		WriteResponse(w, 500, fmt.Sprintf("Can't find instance id: %s", instanceID))
		return
	}

	if instance.Service != service {
		WriteResponse(w, 500, fmt.Sprintf("Wrong ServiceID provided: %q", bindingReq.ServiceID))
		return
	}

	binding := Binding{
		ID:          bindingID,
		AppGUID:     bindingReq.AppGUID,
		PlanID:      bindingReq.PlanID,
		Credentials: `{"password":"letmein"}`,
		Instance:    instance,
	}

	Bindings[bindingID] = &binding
	instance.Bindings = append(instance.Bindings, &binding)

	bindingRes := sbmodel.CreateServiceBindingResponse{
		Credentials: binding.Credentials,
	}

	WriteResponse(w, 200, bindingRes)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v2/catalog", getCatalog).Methods("GET")
	router.HandleFunc("/v2/service_instances/{iid}", createInstance).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{iid}/service_bindings/{bid}", createBinding).Methods("PUT")
	http.Handle("/", router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	fmt.Println("Broker started on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	fmt.Println(err.Error())
}
