package config

import (
	"fmt"
	"sdbi/config"

	"github.com/spf13/cobra"
)

var Delete = &cobra.Command{
	Use:   "delete",
	Short: "config delete",
	Long:  "config delete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.New()
		if err != nil {
			return fmt.Errorf("failed to crarte instance of config: %w", err)
		}

		err = conf.Load()
		if err != nil {
			return fmt.Errorf("faild to load of config.yaml: %w", err)
		}

		configName := args[0]

		err = conf.Delete(configName)
		if err != nil {
			return fmt.Errorf("faield to delete config: %w", err)
		}
		return nil
	},
}
