package cmd

import (
	"fmt"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/fljdin/dispatch/internal/parser"
	"github.com/fljdin/dispatch/internal/tasks"
	"github.com/spf13/cobra"
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

func parseSqlFile(filename string) ([]tasks.Task, error) {
	var finalTasks []tasks.Task

	parser, err := parser.NewBuilder("psql").
		FromFile(filename).
		Build()

	if err != nil {
		return nil, err
	}

	for queryId, query := range parser.Parse() {
		finalTasks = append(finalTasks, tasks.Task{
			QueryID: queryId,
			Name:    fmt.Sprintf("Query loaded from %s", filename),
			Command: tasks.Command{
				Text: query,
				Type: "psql",
				URI:  defaultConnection.CombinedURI(),
			},
		})
	}

	return finalTasks, nil
}

func launch(cmd *cobra.Command, args []string) error {
	config, err := newConfig()
	if err != nil {
		return err
	}

	tasks, err := config.Tasks()
	if err != nil {
		return err
	}

	if len(argSqlFilename) > 0 {
		loadedTasks, err := parseSqlFile(argSqlFilename)
		if err != nil {
			return err
		}
		tasks = append(tasks, loadedTasks...)
	}

	if len(tasks) == 0 {
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

	for _, t := range tasks {
		dispatcher.AddTask(t)
	}

	Debug("loaded tasks =", len(tasks))
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
