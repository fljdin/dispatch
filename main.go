package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ConfigFilename       string
	ConfigFilenameDesc   string = "configuration file"
	configWorkers        int
	ConfigWorkersDefault int    = 2
	ConfigWorkersDesc    string = "number of workers (default 2)"
)

func launch(cmd *cobra.Command, args []string) error {
	configBuild := NewConfigBuilder().FromYAML(ConfigFilename)

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

	dispatcher := NewDispatcher(
		context.Background(),
		config.MaxWorkers,
		len(config.Tasks),
	)

	for i := 0; i < len(config.Tasks); i++ {
		dispatcher.Add(config.Tasks[i])
	}

	dispatcher.Wait()
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
	cmd.Flags().IntVarP(&configWorkers, "jobs", "j", 0, ConfigWorkersDesc)

	// make config flag required by cli
	cmd.Flags().StringVarP(&ConfigFilename, "config", "c", "", ConfigFilenameDesc)
	cmd.MarkFlagRequired("config")

	err := cmd.Execute()
	if err != nil {
		code = 1
	}

	os.Exit(code)
}
