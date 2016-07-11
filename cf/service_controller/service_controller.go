package main

import (
	"flag"
	"fmt"

	"github.com/cncf/servicebroker/cf/service_controller/server"
)

type Options struct {
	ConfigPath string
}

var options Options

func init() {
	flag.StringVar(&options.ConfigPath, "c", ".", "use '-c' option to specify the config file path")

	flag.Parse()
}

func main() {
	s, err := server.CreateServer(server.CreateInMemServiceStorage())
	if err != nil {
		panic(fmt.Sprintf("Error creating server [%s]...", err.Error))
	}

	s.Start()
}
