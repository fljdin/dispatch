package internal

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/lithammer/dedent"
)

var (
	argConfig        string
	argConfigDesc    string = "configuration file"
	argOutput        string
	argOutputDesc    string = "output log file"
	argProcesses     int
	argProcessesDesc string = fmt.Sprintf("number of processes (default %d)", config.ProcessesDefault)
	argVerbose       bool
	argVerboseDesc   string = "verbose mode"
	argVersion       bool
	argVersionDesc   string = "show version"

	usageStr string = dedent.Dedent(`
		Usage:
		  dispatch [options]

		Options:
		  -c, --config=FILE    %s
		  -h, --help           display this help and exit
		  -o, --output=FILE    %s
		  -P, --procs=PROCS    %s
		  -v, --verbose        %s
		      --version        %s
 	`)[1:]

	usage string = fmt.Sprintf(
		usageStr,
		argConfigDesc,
		argOutputDesc,
		argProcessesDesc,
		argVerboseDesc,
		argVersionDesc)
)

func newConfig() (config.Config, error) {
	return config.NewBuilder().
		FromYAML(argConfig).
		WithProcesses(argProcesses).
		WithLogfile(argOutput).
		Build()
}

func Dispatch(version string) {
	parseFlags()
	setupLogging()

	if argVersion {
		fmt.Println(version)
		return
	}

	if argConfig == "" {
		slog.Error("missing configuration file")
		os.Exit(1)
	}

	config, err := newConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	t, err := config.Tasks()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if len(t) == 0 {
		slog.Error("no task to perform")
		os.Exit(1)
	}

	dispatcher, err := dispatcher.NewBuilder().
		WithProcesses(config.Processes).
		WithLogfile(config.Logfile).
		Build()

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	for _, t := range t {
		dispatcher.AddTask(t)
	}

	slog.Debug(
		"loading configuration",
		"tasks", len(t),
		"procs", config.Processes,
	)

	dispatcher.Wait()
	os.Exit(0)
}

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
