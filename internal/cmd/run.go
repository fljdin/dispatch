package cmd

import (
	"context"
	"fmt"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/fljdin/dispatch/internal/models"
	"github.com/fljdin/dispatch/internal/parser"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute tasks",
	RunE:  launch,
}

var defaultConnection models.Connection

func newConfig() (config.Config, error) {
	argPgPassword := ReadHiddenInput("Password: ", argPgPwdPrompt)

	defaultConnection = models.Connection{
		Host:     argPgHost,
		Port:     argPgPort,
		Dbname:   argPgDbname,
		User:     argPgUser,
		Password: argPgPassword,
	}

	return config.NewConfigBuilder().
		FromYAML(argConfigFilename).
		WithMaxWorkers(argMaxWorkers).
		WithLogfile(argLogfile).
		WithDefaultConnection(defaultConnection).
		Build()
}

func newDispatcher(config config.Config) (dispatcher.Dispatcher, error) {
	return dispatcher.NewDispatcherBuilder(context.Background()).
		WithWorkerNumber(config.MaxWorkers).
		WithMemorySize(len(config.Tasks)).
		WithLogfile(config.Logfile).
		WithConsole().
		Build()
}

func parseSqlFile(filename string) ([]models.Task, error) {
	var finalTasks []models.Task

	parser, err := parser.NewParserBuilder("psql").
		FromFile(filename).
		Build()

	if err != nil {
		return nil, err
	}

	for queryId, query := range parser.Parse() {
		finalTasks = append(finalTasks, models.Task{
			QueryID: queryId,
			Type:    "psql",
			Name:    fmt.Sprintf("Query loaded from %s", filename),
			URI:     defaultConnection.CombinedURI(),
			Command: query,
		})
	}

	return finalTasks, nil
}

func launch(cmd *cobra.Command, args []string) error {
	config, err := newConfig()
	if err != nil {
		return err
	}

	if len(argSqlFilename) > 0 {
		tasks, err := parseSqlFile(argSqlFilename)
		if err != nil {
			return err
		}
		config.Tasks = append(config.Tasks, tasks...)
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

	// don't use defaulting feature from cobra as precedence rules are
	// provided by ConfigBuilder
	runCmd.Flags().IntVarP(&argMaxWorkers, "jobs", "j", 0, argMaxWorkersDesc)
	runCmd.Flags().StringVarP(&argSqlFilename, "file", "f", "", argSqlFilenameDesc)
}
