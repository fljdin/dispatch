package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFilename     string
	configFilenameDesc string = "configuration file"
	configWorkers      int
	configWorkersDesc  string = "number of workers (default 2)"
	argVerbose         bool
	argVerboseDesc     string = "verbose mode"
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
	rootCmd.PersistentFlags().StringVarP(&configFilename, "config", "c", "", configFilenameDesc)
	rootCmd.PersistentFlags().BoolVarP(&argVerbose, "verbose", "v", false, argVerboseDesc)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
