package cmd

import (
	"github.com/spf13/cobra"
)

// These are shared between all the commands
var (
	controller string
	timeout    int
	name       string
	spaceGUID  string
)

func init() {
	RootCmd.PersistentFlags().StringVar(&controller, "controller", "http://localhost:10000", "URL for service controller")
	RootCmd.PersistentFlags().IntVar(&timeout, "timeout", 90, "http timeout (in seconds) for interaction with service controller")
}

var RootCmd = &cobra.Command{
	Use:   "sc",
	Short: "CLI for Service Controller operations",
	Long:  "Command Line Interface for the Service Controller",
}
