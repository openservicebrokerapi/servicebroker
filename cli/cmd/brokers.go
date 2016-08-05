package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(brokersCmd)
	brokersCmd.AddCommand(createCmd)
	brokersCmd.AddCommand(describeCmd)
	brokersCmd.AddCommand(listCmd)
	brokersCmd.AddCommand(deleteCmd)
}

var brokersCmd = &cobra.Command{
	Use:   "brokers",
	Short: "manage brokers associated with a service controller",
	Long:  "Manage brokers associated with a service controller",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Add a service broker to service controller",
	Long:  "Add a service broker to service controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("Not implemented yet")
	},
}

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe the specified service broker",
	Long:  "Describe the specified service broker",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("need name of the broker")
		}
		return callService(fmt.Sprintf("/v2/service_brokers/%s", args[0]), "GET", "describe broker", nil)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List brokers the service controller knows about",
	Long:  "List brokers the service controller knows about",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("doing list %s %d\n", controller, timeout)
		return callService("/v2/service_brokers", "GET", "list service brokers", nil)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a broker from the service controller",
	Long:  "Delete a broker from the service controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("doing delete %s %d\n", controller, timeout)
		return nil
	},
}
