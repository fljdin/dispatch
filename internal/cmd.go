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
	argConfigFilename     string
	argConfigFilenameDesc string = "configuration file"
	argOutputFile         string
	argOutputFileDesc     string = "output log file"
	argMaxWorkers         int
	argMaxWorkersDesc     string = "number of workers (default 2)"
	argVerbose            bool
	argVerboseDesc        string = "verbose mode"
	argVersion            bool
	argVersionDesc        string = "show version"

	usageStr string = dedent.Dedent(`
	Usage:
	  dispatch [options]
	
	Options:
	  -j, --jobs <number>        %s
	  -c, --config <filename>    %s
	  -o, --output <filename>    %s
	  -v, --verbose              %s
	  -h, --help                 show this help
	      --version              %s
 	`)[1:]

	usage string = fmt.Sprintf(
		usageStr,
		argMaxWorkersDesc,
		argConfigFilenameDesc, argOutputFileDesc,
		argVerboseDesc, argVersionDesc)
)

func newConfig() (config.Config, error) {
	return config.NewBuilder().
		FromYAML(argConfigFilename).
		WithMaxWorkers(argMaxWorkers).
		WithLogfile(argOutputFile).
		Build()
}

func Dispatch(version string) {
	parseFlags()
	setupLogging()

	if argVersion {
		fmt.Println(version)
		return
	}

	if argConfigFilename == "" {
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
		WithWorkerNumber(config.MaxWorkers).
		WithLogfile(config.Logfile).
		WithConsole().
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
		"workers", config.MaxWorkers,
	)

	dispatcher.Wait()
	os.Exit(0)
}

func parseFlags() {
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}

	flag.BoolVar(&argVersion, "version", false, argVersionDesc)

	flag.IntVar(&argMaxWorkers, "j", 0, argMaxWorkersDesc)
	flag.IntVar(&argMaxWorkers, "jobs", 0, argMaxWorkersDesc)

	flag.StringVar(&argConfigFilename, "c", "", argConfigFilenameDesc)
	flag.StringVar(&argConfigFilename, "config", "", argConfigFilenameDesc)

	flag.StringVar(&argOutputFile, "o", "", argOutputFileDesc)
	flag.StringVar(&argOutputFile, "output", "", argOutputFileDesc)

	flag.BoolVar(&argVerbose, "v", false, argVerboseDesc)
	flag.BoolVar(&argVerbose, "verbose", false, argVerboseDesc)

	flag.Parse()
}
