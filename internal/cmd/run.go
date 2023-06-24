package cmd

import (
	"fmt"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/fljdin/dispatch/internal/tasks"
	"github.com/spf13/cobra"
)

var (
	argMaxWorkers      int
	argMaxWorkersDesc  string = "number of workers (default 2)"
	argSqlFilename     string
	argSqlFilenameDesc string = "file containing SQL statements"
	argType            string
	argTypeDesc        string = "parser type (default sh)"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute tasks",
	RunE:  launch,
}

var defaultConnection tasks.Connection

func newConfig() (config.Config, error) {
	argPgPassword := ReadHiddenInput("Password: ", argPgPwdPrompt)

	defaultConnection = tasks.Connection{
		Host:     argPgHost,
		Port:     argPgPort,
		Dbname:   argPgDbname,
		User:     argPgUser,
		Password: argPgPassword,
	}

	return config.NewBuilder().
		FromYAML(argConfigFilename).
		WithMaxWorkers(argMaxWorkers).
		WithLogfile(argLogfile).
		WithDefaultConnection(defaultConnection).
		Build()
}

func launch(cmd *cobra.Command, args []string) error {
	config, err := newConfig()
	if err != nil {
		return err
	}

	t, err := config.Tasks()
	if err != nil {
		return err
	}

	if len(argSqlFilename) > 0 {
		t = append(t, tasks.Task{
			ID: 1,
			Command: tasks.Command{
				Type: argType,
				File: argSqlFilename,
				URI:  config.DefaultConnection.CombinedURI(),
			},
		})
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

	// don't use defaulting feature from cobra as precedence rules are
	// provided by ConfigBuilder
	runCmd.Flags().IntVarP(&argMaxWorkers, "jobs", "j", 0, argMaxWorkersDesc)
	runCmd.Flags().StringVarP(&argSqlFilename, "file", "f", "", argSqlFilenameDesc)
	runCmd.Flags().StringVarP(&argType, "type", "t", "sh", argTypeDesc)
}
