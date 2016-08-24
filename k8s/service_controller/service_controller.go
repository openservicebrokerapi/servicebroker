package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/cncf/servicebroker/k8s/service_controller/server"
	"github.com/cncf/servicebroker/k8s/service_controller/server/k8s"
	"github.com/cncf/servicebroker/k8s/service_controller/server/mem"
)

type Options struct {
	ConfigPath string
	Backend    string
}

var options Options

func init() {
	flag.StringVar(&options.ConfigPath, "c", ".", "use '-c' option to specify the config file path")
	flag.StringVar(&options.Backend, "backend", "mem", "backend to use for storing info. 'mem' for in memory backend (default), 'k8s' for kubernetes based backend")
	flag.Parse()
}

type Meta struct {
	Name string `json:"name"`
}

type VName struct {
	Name string `json:"name"`
}

type TPR struct {
	Meta       `json:"metadata"`
	ApiVersion string  `json:"apiVersion"`
	kind       string  `json:"kind"`
	Versions   []VName `json:"versions"`
}

const apiVersion string = "extensions/v1beta1"
const thirdPartyResourceString string = "ThirdPartyResource"

var versionMap []VName = []VName{{"v1alpha1"}}
var serviceBrokerMeta Meta = Meta{"service-broker.cncf.org"}
var serviceBindingMeta Meta = Meta{"service-binding.cncf.org"}
var serviceInstanceMeta Meta = Meta{"service-instance.cncf.org"}
var serviceBrokerDefinition TPR = TPR{serviceBrokerMeta, apiVersion, thirdPartyResourceString, versionMap}
var serviceBindingDefinition TPR = TPR{serviceBindingMeta, apiVersion, thirdPartyResourceString, versionMap}
var serviceInstanceDefinition TPR = TPR{serviceInstanceMeta, apiVersion, thirdPartyResourceString, versionMap}

func main() {
	var s *server.Server
	var err error
	switch {
	case options.Backend == "mem":
		s, err = server.CreateServer(mem.CreateInMemServiceStorage())
	case options.Backend == "k8s":
		s, err = server.CreateServer(k8s.CreateServiceStorage())
		// define the resources once at startup
		// results in ServiceBrokers

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(&serviceBrokerDefinition)
		fmt.Printf("encoded bytes: %v\n", b.String())
		r, e := http.Post("http://127.0.0.1:8080/apis/extensions/v1beta1/thirdpartyresources", "application/json", b)
		fmt.Printf("result: %v\n", r)
		if nil != e || 201 != r.StatusCode {
			fmt.Printf("Error creating k8s TPR [%s]...\n%v\n", e, r)
		}
		// serviceBindingDefinition
		// serviceInstanceDefinition

		// cleanup afterwards by `kubectl delete thirdpartyresource service-broker.cncf.org`
	}
	if err != nil {
		panic(fmt.Sprintf("Error creating server [%s]...", err.Error))
	}

	s.Start()
}
