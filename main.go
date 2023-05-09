package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/spf13/cobra"
)

var (
	ConfigFilename     string
	ConfigFilenameDesc string = "configuration file"
	configWorkers      int
	ConfigWorkersDesc  string = "number of workers (default 2)"
)

func launch(cmd *cobra.Command, args []string) error {
	configBuild := config.NewConfigBuilder().FromYAML(ConfigFilename)

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

	dispatcher, err := dispatcher.NewDispatcherBuilder(context.Background()).
		WithWorkerNumber(config.MaxWorkers).
		WithMemorySize(len(config.Tasks)).
		WithTraceFile(config.Summary).
		WithConsole().
		Build()

	if err != nil {
		return err
	}

	for _, t := range config.Tasks {
		dispatcher.AddTask(t)
	}

	dispatcher.Wait()
	return nil
}

func main() {
	cmd := &cobra.Command{
		Use:  "dispatch -c config [-j 2]",
		RunE: launch,
	}

	// don't use defaulting feature from cobra
	// precedence rules are provided by ConfigBuilder
	cmd.Flags().IntVarP(&configWorkers, "jobs", "j", 0, ConfigWorkersDesc)

	// make the config flag required by cli
	cmd.Flags().StringVarP(&ConfigFilename, "config", "c", "", ConfigFilenameDesc)
	cmd.MarkFlagRequired("config")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
