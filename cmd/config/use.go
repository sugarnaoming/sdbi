package config

import (
	"fmt"
	"sdbi/config"

	"github.com/spf13/cobra"
)

var Use = &cobra.Command{
	Use:   "use",
	Short: "config use",
	Long:  "config use",
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

		newConfigName := args[0]

		err = conf.Use(newConfigName)
		if err != nil {
			return fmt.Errorf("failed to select config: %w", err)
		}

		return nil
	},
}
