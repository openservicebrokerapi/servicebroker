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

	// Broker related stuff
	router.HandleFunc("/v2/service_brokers", s.controller.ListServiceBrokers).Methods("GET")
	router.HandleFunc("/v2/service_brokers/{broker_name}/inventory", s.controller.Inventory).Methods("GET")
	router.HandleFunc("/v2/service_brokers/{broker_name}", s.controller.GetServiceBroker).Methods("GET")
	router.HandleFunc("/v2/service_brokers/{broker_name}", s.controller.CreateServiceBroker).Methods("POST")
	router.HandleFunc("/v2/service_brokers/{broker_name}", s.controller.DeleteServiceBroker).Methods("DELETE")

	router.HandleFunc("/v2/service_brokers/{broker_name}/service_instances/", s.controller.ListServiceInstances).Methods("GET")
	router.HandleFunc("/v2/service_brokers/{broker_name}/service_instances/{service_name}", s.controller.GetServiceInstance).Methods("GET")
	router.HandleFunc("/v2/service_brokers/{broker_name}/service_instances/{service_name}", s.controller.CreateServiceInstance).Methods("POST")
	router.HandleFunc("/v2/service_brokers/{broker_name}/service_instances/{service_name}", s.controller.DeleteServiceInstance).Methods("DELETE")

	router.HandleFunc("/v2/service_brokers/{broker_name}/service_instances/{service_name}/service_bindings", s.controller.ListServiceBindings).Methods("GET")
	router.HandleFunc("/v2/service_brokers/{broker_name}/service_instances/{service_name}/service_bindings/{service_binding_guid}", s.controller.GetServiceBinding).Methods("GET")
	router.HandleFunc("/v2/service_brokers/{broker_name}/service_instances/{service_name}/service_bindings/{service_binding_guid}", s.controller.CreateServiceBinding).Methods("POST")
	router.HandleFunc("/v2/service_brokers/{broker_name}/service_instances/{service_name}/service_bindings/{service_binding_guid}", s.controller.DeleteServiceBinding).Methods("DELETE")

	http.Handle("/", router)

	cfPort := os.Getenv("PORT")
	if cfPort == "" {
		cfPort = "10000"
	}
	fmt.Println("Server started on port " + cfPort)
	err := http.ListenAndServe(":"+cfPort, nil)
	fmt.Println(err.Error())
}
