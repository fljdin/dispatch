package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configFilename string
var configFilenameDesc string = "configuration file"

func launch(cmd *cobra.Command, args []string) error {
	config, err := NewConfigBuilder().FromYAML(configFilename).Build()
	if err != nil {
		return err
	}

	fmt.Println("Config loaded with", len(config.Tasks), "tasks")
	return nil
}

func main() {
	code := 0
	cmd := &cobra.Command{
		Use:  "-c config",
		RunE: launch,
	}

	cmd.Flags().StringVarP(&configFilename, "config", "c", "", configFilenameDesc)
	cmd.MarkFlagRequired("config")

	err := cmd.Execute()
	if err != nil {
		code = 1
	}

	os.Exit(code)
}
