package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	SERVICE_PLANS_URL = "/v2/service_plans"
)

func init() {
	RootCmd.AddCommand(servicePlansCmd)
	servicePlansCmd.AddCommand(listServicePlansCmd)
	servicePlansCmd.AddCommand(describeServicePlansCmd)
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

var describeServicePlansCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a service plan",
	Long:  "Describe a service plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("Need <service> <plan>")
		}

		inv, err := fetchInventory()
		if err != nil {
			return err
		}

		for _, s := range inv.Services {
			if s.Name == args[0] {
				for _, p := range s.Plans {
					if p.Name == args[1] {
						fmt.Printf("Schema:\n%s", p.Schemas.Instance)
						return nil
					}
				}
			}
		}
		return fmt.Errorf("Can't find a service / plan : %s/%s", args[0], args[1])

		return callService(SERVICE_PLANS_URL, "GET", "list service plans", nil)
	},
}
