package main

import (
	"flag"
	"fmt"

	"github.com/cncf/servicebroker/k8s/service_controller/server"
	"github.com/cncf/servicebroker/k8s/service_controller/server/k8s"
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
		s, err = server.CreateServer(server.CreateInMemServiceStorage())
	case options.Backend == "k8s":
		s, err = server.CreateServer(k8s.CreateServiceStorage())
	}
	if err != nil {
		panic(fmt.Sprintf("Error creating server [%s]...", err.Error))
	}

	s.Start()
}
