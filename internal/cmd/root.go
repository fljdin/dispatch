package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	argConfigFilename     string
	argConfigFilenameDesc string = "configuration file"
	argLogfile            string
	argLogfileDesc        string = "log file"
	argMaxWorkers         int
	argMaxWorkersDesc     string = "number of workers (default 2)"
	argPgDbname           string
	argPgDbnameDesc       string = "database name to connect to"
	argPgHost             string
	argPgHostDesc         string = "database server host or socket directory"
	argPgPort             int
	argPgPortDesc         string = "database server port"
	argPgPwdPrompt        bool
	argPgPwdPromptdDesc   string = "force password prompt"
	argPgUser             string
	argPgUserDesc         string = "database user name"
	argSqlFilename        string
	argSqlFilenameDesc    string = "file containing SQL statements"
	argVerbose            bool
	argVerboseDesc        string = "verbose mode"
)

var RootCmd = &cobra.Command{
	Use:   "dispatch",
	Short: "Dispatch tasks described in a YAML file",
}

func Debug(data ...any) {
	if argVerbose {
		data = append([]any{"DEBUG"}, data...)
		log.Println(data...)
	}
}

func ReadInput(prompt string, condition bool) string {
	var value string
	if condition {
		fmt.Print(prompt)

		reader := bufio.NewReader(os.Stdin)
		value, _ = reader.ReadString('\n')
		value = strings.TrimSpace(value)
	}
	return value
}

func ReadHiddenInput(prompt string, condition bool) string {
	var value string
	if condition {
		fmt.Print(prompt)

		reader, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return value
		}
		value = string(reader)
		value = strings.TrimSpace(value)
		fmt.Print("\n")
	}
	return value
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&argConfigFilename, "config", "c", "", argConfigFilenameDesc)
	RootCmd.PersistentFlags().BoolVarP(&argVerbose, "verbose", "v", false, argVerboseDesc)
	RootCmd.PersistentFlags().StringVarP(&argLogfile, "log", "l", "", argLogfileDesc)

	RootCmd.PersistentFlags().StringVarP(&argPgHost, "host", "h", "", argPgHostDesc)
	RootCmd.PersistentFlags().IntVarP(&argPgPort, "port", "p", 0, argPgPortDesc)
	RootCmd.PersistentFlags().StringVarP(&argPgDbname, "dbname", "d", "", argPgDbnameDesc)
	RootCmd.PersistentFlags().StringVarP(&argPgUser, "user", "U", "", argPgUserDesc)
	RootCmd.PersistentFlags().BoolVarP(&argPgPwdPrompt, "password", "W", false, argPgPwdPromptdDesc)

	RootCmd.PersistentFlags().Bool("help", false, "show help")
	RootCmd.CompletionOptions.DisableDefaultCmd = true
}
