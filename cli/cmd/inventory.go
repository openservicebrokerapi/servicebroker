package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const (
	INVENTORY_FMT_STR = "%-20s %-25s %-30s\n"
)

func init() {
	RootCmd.AddCommand(inventoryCmd)
}

var inventoryCmd = &cobra.Command{
	Aliases: []string{"inv"},
	Use:     "inventory",
	Short:   "List the available services",
	Long:    "List the available services",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := fetchInventory()
		if err != nil {
			return err
		}
		fmt.Printf(INVENTORY_FMT_STR, "service", "plans", "description")
		for _, s := range c.Services {
			var plans []string
			for _, p := range s.Plans {
				plans = append(plans, p.Name)
			}
			fmt.Printf(INVENTORY_FMT_STR, s.Name, strings.Join(plans, ","), s.Description)
		}
		return nil
	},
}
