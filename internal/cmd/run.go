package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/fljdin/dispatch/internal/models"
	"github.com/fljdin/dispatch/internal/parser"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute tasks",
	RunE:  run,
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

func newReader() (io.Reader, error) {
	if len(argSqlFilename) > 0 && argSqlFilename != "-" {
		file, err := os.Open(argSqlFilename)
		if err != nil {
			return nil, fmt.Errorf("failed open file: %v", err)
		}
		return file, nil
	}
	return nil, nil
}

func parseInput(input io.Reader) ([]models.Task, error) {
	var finalTasks []models.Task

	content, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %v", err)
	}

	parser, _ := parser.NewParserBuilder("psql").
		WithContent(string(content)).
		Build()

	for queryId, query := range parser.Parse() {
		finalTasks = append(finalTasks, models.Task{
			QueryID: queryId,
			Name:    "Query loaded from input",
			Command: models.Command{
				Text: query,
				Type: "psql",
				URI:  defaultConnection.CombinedURI(),
			},
		})
	}

	return finalTasks, nil
}

func run(cmd *cobra.Command, args []string) error {
	config, err := newConfig()
	if err != nil {
		return err
	}

	tasks, err := config.GetTasks()
	if err != nil {
		return err
	}

	inputReader, err := newReader()
	if err != nil {
		return err
	}

	if inputReader == nil && HasInput() {
		inputReader = cmd.InOrStdin()
	}

	if inputReader != nil {
		loadedTasks, err := parseInput(inputReader)

		if err != nil {
			return err
		}

		tasks = append(tasks, loadedTasks...)
	}

	if len(tasks) == 0 {
		return fmt.Errorf("no task to perform")
	}

	dispatcher, err := dispatcher.NewDispatcherBuilder(context.Background()).
		WithWorkerNumber(config.MaxWorkers).
		WithMemorySize(len(tasks)).
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
