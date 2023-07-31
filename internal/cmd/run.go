package cmd

import (
	"fmt"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tasks from configuration file",
	RunE:  run,
}

func newConfig() (config.Config, error) {
	return config.NewBuilder().
		FromYAML(argConfigFilename).
		WithMaxWorkers(argMaxWorkers).
		WithLogfile(argLogfile).
		WithDefaultConnection(DefaultConnection()).
		Build()
}

func run(cmd *cobra.Command, args []string) error {
	config, err := newConfig()
	if err != nil {
		return err
	}

	t, err := config.Tasks()
	if err != nil {
		return err
	}

	if len(t) == 0 {
		return fmt.Errorf("no task to perform")
	}

	dispatcher, err := dispatcher.NewBuilder().
		WithWorkerNumber(config.MaxWorkers).
		WithLogfile(config.Logfile).
		WithConsole().
		Build()

	if err != nil {
		return err
	}

	for _, t := range t {
		dispatcher.AddTask(t)
	}

	Debug("loaded tasks =", len(t))
	Debug("max workers =", config.MaxWorkers)

	dispatcher.Wait()
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&argParserType, "type", "t", "sh", argParserTypeDesc)
}
