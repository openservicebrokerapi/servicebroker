package main

import (
	"flag"
	"fmt"

	"github.com/servicebroker/servicebroker/k8s/service_controller/server"
	"github.com/servicebroker/servicebroker/k8s/service_controller/server/k8s"
	"github.com/servicebroker/servicebroker/k8s/service_controller/server/mem"
)

type Options struct {
	ConfigPath string
	Backend    string
	Host       string
}

var options Options

func init() {
	flag.StringVar(&options.ConfigPath, "c", ".", "use '-c' option to specify the config file path")
	flag.StringVar(&options.Backend, "backend", "mem", "backend to use for storing info. 'mem' for in memory backend (default), 'k8s' for kubernetes based backend")
	flag.StringVar(&options.Host, "host", "localhost:8080", "the host of the backend to connect to")
	flag.Parse()
}

func main() {
	var s *server.Server
	var err error
	switch {
	case options.Backend == "mem":
		s, err = server.CreateServer(mem.CreateInMemServiceStorage())
	case options.Backend == "k8s":
		s, err = server.CreateServer(k8s.CreateServiceStorage(options.Host))
	}
	if err != nil {
		panic(fmt.Sprintf("Error creating server [%s]...", err.Error))
	}

	s.Start()
}
