package internal

import (
	"flag"
	"fmt"

	"github.com/fljdin/dispatch/internal/config"
)

var (
	argConfig        string
	argConfigDesc    string = "configuration file"
	argOutput        string
	argOutputDesc    string = "redirect output to file"
	argProcesses     int
	argProcessesDesc string = fmt.Sprintf("number of processes (default %d)", config.ProcessesDefault)
	argVerbose       bool
	argVerboseDesc   string = "verbose mode"
	argVersion       bool
	argVersionDesc   string = "show version"
)

func parseFlags() {
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}

	flag.StringVar(&argConfig, "c", "", argConfigDesc)
	flag.StringVar(&argConfig, "config", "", argConfigDesc)

	flag.StringVar(&argOutput, "o", "", argOutputDesc)
	flag.StringVar(&argOutput, "output", "", argOutputDesc)

	flag.IntVar(&argProcesses, "P", 0, argProcessesDesc)
	flag.IntVar(&argProcesses, "procs", 0, argProcessesDesc)

	flag.BoolVar(&argVerbose, "v", false, argVerboseDesc)
	flag.BoolVar(&argVerbose, "verbose", false, argVerboseDesc)

	flag.BoolVar(&argVersion, "version", false, argVersionDesc)

	flag.Parse()
}

func newConfig() (config.Config, error) {
	return config.NewBuilder().
		FromYAML(argConfig).
		WithProcesses(argProcesses).
		WithLogfile(argOutput).
		WithVerbose(argVerbose).
		Build()
}
