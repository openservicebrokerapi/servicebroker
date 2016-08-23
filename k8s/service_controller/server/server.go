package server

import (
	//	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Server struct {
	controller *Controller
}

func CreateServer(serviceStorage ServiceStorage) (*Server, error) {
	return &Server{
		controller: CreateController(serviceStorage),
	}, nil
}

func (s *Server) Start() {
	router := mux.NewRouter()

	// TODO: the actual inventory API should be /v2/services[/...] and
	// /v2/service_plans[/...].
	router.HandleFunc("/v2/service_plans", s.controller.Inventory).Methods("GET")

	// Broker related stuff
	router.HandleFunc("/v2/service_brokers", s.controller.ListServiceBrokers).Methods("GET")
	router.HandleFunc("/v2/service_brokers", s.controller.CreateServiceBroker).Methods("POST")
	router.HandleFunc("/v2/service_brokers/{broker_id}", s.controller.GetServiceBroker).Methods("GET")
	router.HandleFunc("/v2/service_brokers/{broker_id}", s.controller.DeleteServiceBroker).Methods("DELETE")
	// TODO: implement updating a service broker.
	// router.HandleFunc("/v2/service_brokers/{broker_id}", s.Controller.UpdateServiceBroker).Methods.("PUT")

	router.HandleFunc("/v2/service_instances", s.controller.ListServiceInstances).Methods("GET")
	router.HandleFunc("/v2/service_instances", s.controller.CreateServiceInstance).Methods("POST")
	router.HandleFunc("/v2/service_instances/{service_id}", s.controller.GetServiceInstance).Methods("GET")
	router.HandleFunc("/v2/service_instances/{service_id}", s.controller.DeleteServiceInstance).Methods("DELETE")
	// TODO: implement list service bindings for this service instance.
	// router.HandleFunc("/v2/service_instances/{service_id}/service_bindings", s.controller.ListServiceInstanceBindings).Methods("GET")

	router.HandleFunc("/v2service_bindings", s.controller.ListServiceBindings).Methods("GET")
	router.HandleFunc("/v2/service_bindings", s.controller.CreateServiceBinding).Methods("POST")
	router.HandleFunc("/v2service_bindings/{binding_id}", s.controller.GetServiceBinding).Methods("GET")
	router.HandleFunc("/v2/service_bindings/{binding_id}", s.controller.DeleteServiceBinding).Methods("DELETE")

	http.Handle("/", router)

	cfPort := os.Getenv("PORT")
	if cfPort == "" {
		cfPort = "10000"
	}
	fmt.Println("Service Controller started on port " + cfPort)
	err := http.ListenAndServe(":"+cfPort, nil)
	fmt.Println(err.Error())
}
