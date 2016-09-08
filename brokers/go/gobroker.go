package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/servicebroker/servicebroker/k8s/service_controller/model"
)

type Service struct {
	Name        string               `json:"name"`
	ID          string               `json:"id"`
	Description string               `json:"description"`
	Bindable    bool                 `json:"bindable"`
	Updateable  bool                 `json:"plan_updateable"`
	Plans       []*model.ServicePlan `json:"plans"`
}

type Plan struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"description"`
	Free        bool   `json:"free"`
}

type Instance struct {
	ID           string
	Service      *Service
	DashboardURL string
	State        string
}

var Services = map[string]*Service{}
var Instances = map[string]*Instance{}

func init() {
	plan1 := model.ServicePlan{
		Name:        "freePlan",
		ID:          "ffffffff-0000-0000-0000-000000000000",
		Description: "free is good",
		Free:        true,
	}
	plan2 := model.ServicePlan{
		Name:        "costlyPlan",
		ID:          "eeeeeeee-1111-1111-1111-111111111111",
		Description: "not so free",
	}

	Service := Service{
		Name:        "myService1",
		ID:          "12345678-abcd-1234-bcde-1234567890ab",
		Description: "something cool",
		Bindable:    true,
		Updateable:  false,
		Plans:       []*model.ServicePlan{&plan1, &plan2},
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
	fmt.Fprintf(w, string(data)+"\n")
}

func getCatalog(w http.ResponseWriter, r *http.Request) {
	// GET /v2/catalog
	Log(r.Method + ": " + r.URL.String())
	data, err := json.Marshal(Services)
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

	svcID := mux.Vars(r)["id"]

	service, ok := Services[svcID]
	if !ok {
		WriteResponse(w, 500, fmt.Sprintf("Can't find service id: %s", svcID))
		return
	}

	var instanceReq model.CreateServiceInstanceRequest
	err := json.NewDecoder(r.Body).Decode(&instanceReq)
	if err != nil {
		WriteResponse(w, 500, err)
		return
	}

	instance := Instance{
		ID:           "123",
		Service:      service,
		DashboardURL: "http://example.com/dashAwayAll",
		State:        "created",
	}

	Instances[instance.ID] = &instance

	instanceRes := model.CreateServiceInstanceResponse{
		DashboardURL:  instance.DashboardURL,
		LastOperation: nil,
	}

	WriteResponse(w, 200, instanceRes)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v2/catalog", getCatalog).Methods("GET")
	router.HandleFunc("/v2/service_instance/{id}", createInstance).Methods("PUT")
	http.Handle("/", router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	fmt.Println("Broker started on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	fmt.Println(err.Error())
}
