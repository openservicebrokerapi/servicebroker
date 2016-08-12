package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	version = "0.0.1"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of sc-cli",
	Long:  "Print the version of sc-cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", version)
	},
}
