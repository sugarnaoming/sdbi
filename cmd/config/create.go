package config

import (
	"fmt"
	"sdbi/config"

	"github.com/spf13/cobra"
)

var Create = &cobra.Command{
	Use:   "create",
	Short: "config create",
	Long:  "config create {config name} {API URL} {token} {UI URL}",
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.New()
		if err != nil {
			return fmt.Errorf("failed to crarte instance of config: %w", err)
		}

		err = conf.Load()
		if err != nil {
			return fmt.Errorf("faild to load of config.yaml: %w", err)
		}

		bp := config.Blueprint{
			ConfigName: args[0],
			APIURL:     args[1],
			Token:      args[2],
			UIURL:      args[3],
		}

		err = conf.Create(bp)
		if err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}

		return nil
	},
}
