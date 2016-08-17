package cmd

import (
	"github.com/spf13/cobra"
)

const (
	SERVICE_PLANS_URL = "/v2/service_plans"
)

func init() {
	RootCmd.AddCommand(servicePlansCmd)
	servicePlansCmd.AddCommand(listServicePlansCmd)
}

var servicePlansCmd = &cobra.Command{
	Aliases: []string{"sp"},
	Use:     "service-plans",
	Short:   "manage service plans associated with a service controller",
	Long:    "Manage service plans associated with a service controller",
}

var listServicePlansCmd = &cobra.Command{
	Use:   "list",
	Short: "List all service plans available from service controller",
	Long:  "List all service plans available from service controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		return callService(SERVICE_PLANS_URL, "GET", "list service plans", nil)
	},
}
