package internal

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/lithammer/dedent"
)

var (
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

func Dispatch(version string) {
	parseFlags()
	setEnvirons()
	setupLogging(os.Stderr, false)

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

	setupLogging(os.Stderr, config.Verbose)

	t, err := config.Tasks()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if len(t) == 0 {
		slog.Error("no task to perform")
		os.Exit(1)
	}

	if config.Logfile != "" {
		f, err := openOutputFile(config.Logfile)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		setupLogging(f, config.Verbose)
	}

	dispatcher := dispatcher.New(config.Processes)

	for _, t := range t {
		dispatcher.AddTask(t)
	}

	slog.Info(
		"loading configuration",
		"tasks", len(t),
		"procs", config.Processes,
		"verbose", config.Verbose,
	)

	dispatcher.Wait()
	os.Exit(0)
}

func setEnvirons() {
	os.Setenv("PGAPPNAME", "dispatch")
}
