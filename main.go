package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configFilename string
var configFilenameDesc string = "configuration file"
var configWorkers int
var configWorkersDefault int = 2
var configWorkersDesc string = "number of workers (default 2)"

func launch(cmd *cobra.Command, args []string) error {
	configBuild := NewConfigBuilder().FromYAML(configFilename)

	if configWorkers > 0 {
		configBuild = configBuild.
			WithMaxWorkers(configWorkers)
	}

	config, err := configBuild.Build()
	if err != nil {
		return err
	}

	fmt.Println("Config loaded with", len(config.Tasks), "tasks")
	fmt.Println("- max workers =", config.MaxWorkers)

	return nil
}

func main() {
	code := 0
	cmd := &cobra.Command{
		Use:  "dispatch -c config [-j 2]",
		RunE: launch,
	}

	// don't use defaulting feature from cobra
	// precedence rules are provided by ConfigBuilder
	cmd.Flags().IntVarP(&configWorkers, "jobs", "j", 0, configWorkersDesc)

	// make config flag required by cli
	cmd.Flags().StringVarP(&configFilename, "config", "c", "", configFilenameDesc)
	cmd.MarkFlagRequired("config")

	err := cmd.Execute()
	if err != nil {
		code = 1
	}

	os.Exit(code)
}
