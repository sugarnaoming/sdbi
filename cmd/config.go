package cmd

import (
	"fmt"
	"sdbi/cmd/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config commands",
	Long:  "config commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("this command is config")
		return nil
	},
}

func init() {
	configCmd.AddCommand(config.Create)
	configCmd.AddCommand(config.Delete)
	configCmd.AddCommand(config.Use)
	configCmd.AddCommand(config.Set)
	configCmd.AddCommand(config.View)
}
