package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	argConfigFilename     string
	argConfigFilenameDesc string = "configuration file"
	argMaxWorkers         int
	argMaxWorkersDesc     string = "number of workers (default 2)"
	argVerbose            bool
	argVerboseDesc        string = "verbose mode"
	argSqlFilename        string
	argSqlFilenameDesc    string = "file containing SQL statements"
)

var rootCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Dispatch tasks described in a YAML file",
}

func Debug(data ...any) {
	if argVerbose {
		data = append([]any{"DEBUG"}, data...)
		log.Println(data...)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&argConfigFilename, "config", "c", "", argConfigFilenameDesc)
	rootCmd.PersistentFlags().BoolVarP(&argVerbose, "verbose", "v", false, argVerboseDesc)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
