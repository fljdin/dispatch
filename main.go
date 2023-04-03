package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configFile string
var configFileDesc string = "configuration file"

func launch(cmd *cobra.Command, args []string) {
	config, err := NewConfigBuilder().FromYAML(configFile).Build()

	fmt.Println("Config loaded with % tasks", len(config.Tasks))
}

func execute() int {
	code := 0
	cmd := &cobra.Command{
		Use: "-c config",
		Run: launch,
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "", configFileDesc)
	cmd.MarkFlagRequired("config")

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		code = 1
	}

	return code
}

func main() {
	os.Exit(execute())
}
