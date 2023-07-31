package cmd

import (
	"fmt"

	"github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/fljdin/dispatch/internal/tasks"
	"github.com/fljdin/dispatch/internal/tasks/actions"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute tasks from commands",
	RunE:  exec,
}

func exec(cmd *cobra.Command, args []string) error {
	t := []tasks.Task{}
	defaultConnection := DefaultConnection()

	if len(argFilename) > 0 {
		t = append(t, tasks.Task{
			ID: 1,
			Action: actions.FileLoader{
				File: argFilename,
				Type: argParserType,
				URI:  defaultConnection.CombinedURI(),
			},
		})
	} else if len(argCommand) > 0 {
		t = append(t, tasks.Task{
			ID: 1,
			Action: actions.OutputLoader{
				Text: argCommand,
				From: argParserType,
				Type: argExecType,
				URI:  defaultConnection.CombinedURI(),
			},
		})
	}

	if len(t) == 0 {
		return fmt.Errorf("no task to perform")
	}

	dispatcher, err := dispatcher.NewBuilder().
		WithWorkerNumber(argMaxWorkers).
		WithLogfile(argLogfile).
		WithConsole().
		Build()

	if err != nil {
		return err
	}

	for _, t := range t {
		dispatcher.AddTask(t)
	}

	Debug("loaded tasks =", len(t))
	Debug("max workers =", argMaxWorkers)

	dispatcher.Wait()
	return nil
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().StringVarP(&argCommand, "command", "C", "", argCommandDesc)
	execCmd.Flags().StringVarP(&argFilename, "file", "f", "", argFilenameDesc)
	execCmd.Flags().StringVarP(&argExecType, "to", "T", "sh", argExecTypeDesc)
	execCmd.Flags().StringVarP(&argParserType, "type", "t", "sh", argParserTypeDesc)
}
