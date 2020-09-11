package config

import (
	"fmt"
	"sdbi/config"

	"github.com/spf13/cobra"
)

var ApiURL string
var token string

var Set = &cobra.Command{
	Use:     "set",
	Short:   "config set {configName}",
	Long:    "config set",
	Example: "config set TestConfig -a http://exmaple.com -t your-token",
	Args:    cobra.ExactArgs(1),
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
			APIURL:     ApiURL,
			Token:      token,
		}

		err = conf.Set(bp)
		if err != nil {
			return fmt.Errorf("failed to set config: %w", err)
		}

		return nil
	},
}

func init() {
	Set.Flags().StringVarP(&ApiURL, "api-url", "a", "", "API URL")
	Set.Flags().StringVarP(&token, "token", "t", "", "Token")
}
