package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sdbi",
	Short: "sdbi root",
	Long:  "sdbi root",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("This command is sdbi")
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add this below in function when new command added
	// e.g. rootCmd.AddCommand(COMMAND_NAME)
	rootCmd.AddCommand(configCmd)
}
