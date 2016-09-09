package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	sbmodel "github.com/servicebroker/servicebroker/model/service_broker"
	"github.com/spf13/cobra"
)

const (
	BROKERS_URL     = "/v2/service_brokers"
	BROKERS_FMT_STR = "/v2/service_brokers/%s"
)

var (
	user     string
	password string
)

func init() {
	RootCmd.AddCommand(brokersCmd)

	brokersCmd.AddCommand(createBrokersCmd)
	brokersCmd.AddCommand(describeBrokersCmd)
	brokersCmd.AddCommand(listBrokersCmd)
	brokersCmd.AddCommand(deleteBrokersCmd)

	createBrokersCmd.Flags().StringVarP(&user, "user", "u", "", "user name to use for broker auth")
	createBrokersCmd.Flags().StringVarP(&password, "password", "p", "", "password to use for broker auth")
	createBrokersCmd.Flags().StringVarP(&spaceGUID, "spaceGUID", "s", "default", "Space GUID to use for broker")

}

var brokersCmd = &cobra.Command{
	Use:   "brokers",
	Short: "manage brokers associated with a service controller",
	Long:  "Manage brokers associated with a service controller",
}

var createBrokersCmd = &cobra.Command{
	Use:   "create <NAME> <BROKER_URL>",
	Short: "Add a service broker to service controller",
	Long:  "Add a service broker to service controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("need NAME and BROKER_URL")
		}
		req := sbmodel.CreateServiceBrokerRequest{
			Name:      args[0],
			BrokerURL: args[1],
		}
		if len(user) > 0 {
			req.AuthUsername = user
		}
		if len(password) > 0 {
			req.AuthPassword = password
		}
		if len(spaceGUID) > 0 {
			req.SpaceGUID = spaceGUID
		}
		body, err := json.Marshal(req)
		if err != nil {
			return err
		}
		return callService(BROKERS_URL, "POST", "create broker", ioutil.NopCloser(bytes.NewReader(body)))
	},
}

var describeBrokersCmd = &cobra.Command{
	Use:   "describe <NAME>",
	Short: "Describe the specified service broker",
	Long:  "Describe the specified service broker",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("need NAME of the broker")
		}
		guid, err := fetchBrokerGUID(args[0])
		if err != nil {
			return err
		}
		return callService(fmt.Sprintf(BROKERS_FMT_STR, guid), "GET", "describe broker", nil)
	},
}

var listBrokersCmd = &cobra.Command{
	Use:   "list",
	Short: "List brokers the service controller knows about",
	Long:  "List brokers the service controller knows about",
	RunE: func(cmd *cobra.Command, args []string) error {
		return callService(BROKERS_URL, "GET", "list service brokers", nil)
	},
}

var deleteBrokersCmd = &cobra.Command{
	Use:   "delete <NAME>",
	Short: "Delete a broker from the service controller",
	Long:  "Delete a broker from the service controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("need NAME of the broker to delete")
		}

		guid, err := fetchBrokerGUID(args[0])
		if err != nil {
			return err
		}
		return callService(fmt.Sprintf(BROKERS_FMT_STR, guid), "DELETE", "delete broker", nil)
	},
}
