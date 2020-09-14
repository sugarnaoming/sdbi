package config

import (
	"fmt"
	"sdbi/config"
	"sdbi/display"

	"github.com/spf13/cobra"
)

var View = &cobra.Command{
	Use:   "view",
	Short: "config view",
	Long:  "config view",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := config.New()
		if err != nil {
			return fmt.Errorf("failed to crarte instance of config: %w", err)
		}

		err = conf.Load()
		if err != nil {
			return fmt.Errorf("faild to load of config.yaml: %w", err)
		}

		d := display.New()
		d.CreateHeader("Current", "ConfigName", "Setting")
		rowContents := []display.RowContents{}
		for k, v := range conf.Config {
			var contents display.RowContents
			if k == conf.CurrentConfName {
				// Primary setting
				contents = d.CreateRowContents("*", k, fmt.Sprintf("API URL:\t%s", v.APIURL))
			} else {
				// Not primary setting
				contents = d.CreateRowContents("", k, fmt.Sprintf("API URL:\t%s", v.APIURL))
			}
			rowContents = append(rowContents, contents)
			contents = d.CreateRowContents("", "", fmt.Sprintf("User Token:\t%s", v.UserToken))
			rowContents = append(rowContents, contents)
			contents = d.CreateRowContents("", "", fmt.Sprintf("UI URL:\t%s", v.UIURL))
			rowContents = append(rowContents, contents)
		}

		err = d.View(rowContents)
		if err != nil {
			return fmt.Errorf("faled to display: %w", err)
		}
		return nil
	},
}
