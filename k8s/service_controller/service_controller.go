package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/cncf/servicebroker/k8s/service_controller/server"
	"github.com/cncf/servicebroker/k8s/service_controller/server/k8s"
	"github.com/cncf/servicebroker/k8s/service_controller/server/mem"
	"github.com/cncf/servicebroker/k8s/service_controller/utils"
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
		servicebrokerTPRDefYaml := `metadata:
  name: service-broker.cncf.org
apiVersion: extensions/v1beta1
kind: ThirdPartyResource
versions:
  - name: v1alpha1`

		// 		// results in ServiceInstances
		// 		serviceinstanceTPRDefYaml := `metadata:
		//   name: service-instance.cncf.org
		// apiVersion: extensions/v1beta1
		// kind: ThirdPartyResource
		// versions:
		// - name: v1alpha1`
		// 		// results in SericeBindings
		// 		servicebindingTPRDefYaml := `metadata:
		//   name: service-binding.cncf.org
		// apiVersion: extensions/v1beta1
		// kind: ThirdPartyResource
		// versions:
		// - name: v1alpha1`

		s, e := utils.KubeCreateResource(strings.NewReader(servicebrokerTPRDefYaml))
		if nil != e {
			panic(fmt.Sprintf("Error creating k8s TPR [%s]...\n%v", e, s))
		}
		// s, e = utils.KubeCreateResource(strings.NewReader(serviceinstanceTPRDefYaml))
		// s, e = utils.KubeCreateResource(strings.NewReader(servicebindingTPRDefYaml))

		// cleanup afterwards by `kubectl delete thirdpartyresource service-broker.cncf.org`
	}
	if err != nil {
		panic(fmt.Sprintf("Error creating server [%s]...", err.Error))
	}

	s.Start()
}
