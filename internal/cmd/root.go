package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fljdin/dispatch/internal/tasks"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	argCommand            string
	argCommandDesc        string = "loader command to execute"
	argConfigFilename     string
	argConfigFilenameDesc string = "configuration file"
	argExecType           string
	argExecTypeDesc       string = "executor type"
	argFilename           string
	argFilenameDesc       string = "file containing commands to execute"
	argLogfile            string
	argLogfileDesc        string = "log file"
	argMaxWorkers         int
	argMaxWorkersDesc     string = "number of workers (default 2)"
	argParserType         string
	argParserTypeDesc     string = "parser type"
	argPgDbname           string
	argPgDbnameDesc       string = "database name to connect to"
	argPgHost             string
	argPgHostDesc         string = "database server host or socket directory"
	argPgPort             int
	argPgPortDesc         string = "database server port"
	argPgPwdPrompt        bool
	argPgPwdPromptdDesc   string = "force password prompt"
	argPgService          string
	argPgServiceDesc      string = "database service"
	argPgUser             string
	argPgUserDesc         string = "database user name"
	argVerbose            bool
	argVerboseDesc        string = "verbose mode"
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

func DefaultConnection() tasks.Connection {
	argPgPassword := ReadHiddenInput("Password: ", argPgPwdPrompt)

	return tasks.Connection{
		Service:  argPgService,
		Host:     argPgHost,
		Port:     argPgPort,
		Dbname:   argPgDbname,
		User:     argPgUser,
		Password: argPgPassword,
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&argMaxWorkers, "jobs", "j", 0, argMaxWorkersDesc)
	rootCmd.PersistentFlags().StringVarP(&argConfigFilename, "config", "c", "", argConfigFilenameDesc)
	rootCmd.PersistentFlags().BoolVarP(&argVerbose, "verbose", "v", false, argVerboseDesc)
	rootCmd.PersistentFlags().StringVarP(&argLogfile, "log", "l", "", argLogfileDesc)

	rootCmd.PersistentFlags().StringVarP(&argPgService, "service", "s", "", argPgServiceDesc)
	rootCmd.PersistentFlags().StringVarP(&argPgHost, "host", "h", "", argPgHostDesc)
	rootCmd.PersistentFlags().IntVarP(&argPgPort, "port", "p", 0, argPgPortDesc)
	rootCmd.PersistentFlags().StringVarP(&argPgDbname, "dbname", "d", "", argPgDbnameDesc)
	rootCmd.PersistentFlags().StringVarP(&argPgUser, "user", "U", "", argPgUserDesc)
	rootCmd.PersistentFlags().BoolVarP(&argPgPwdPrompt, "password", "W", false, argPgPwdPromptdDesc)

	rootCmd.PersistentFlags().Bool("help", false, "show help")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
