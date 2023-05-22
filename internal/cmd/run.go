package cmd

import (
	"context"
	"fmt"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute tasks",
	RunE:  launch,
}

func newConfig() (config.Config, error) {
	configBuild := config.NewConfigBuilder()

	if len(configFilename) > 0 {
		configBuild.FromYAML(configFilename)
	}

	if configWorkers > 0 {
		configBuild = configBuild.
			WithMaxWorkers(configWorkers)
	}

	return configBuild.Build()
}

func newDispatcher(config config.Config) (dispatcher.Dispatcher, error) {
	return dispatcher.NewDispatcherBuilder(context.Background()).
		WithWorkerNumber(config.MaxWorkers).
		WithMemorySize(len(config.Tasks)).
		WithTraceFile(config.Summary).
		WithConsole().
		Build()
}

func launch(cmd *cobra.Command, args []string) error {
	config, err := newConfig()
	if err != nil {
		return err
	}

	if len(config.Tasks) == 0 {
		return fmt.Errorf("no task to perform")
	}

	dispatcher, err := newDispatcher(config)
	if err != nil {
		return err
	}

	for _, t := range config.Tasks {
		dispatcher.AddTask(t)
	}

	Debug("loaded tasks =", len(config.Tasks))
	Debug("max workers =", config.MaxWorkers)

	dispatcher.Wait()
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)

	// don't use defaulting feature from cobra
	// precedence rules are provided by ConfigBuilder
	runCmd.Flags().IntVarP(&configWorkers, "jobs", "j", 0, configWorkersDesc)
}
