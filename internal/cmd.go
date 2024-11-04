package internal

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/routines"
	"github.com/knadh/koanf/v2"
	"github.com/lithammer/dedent"
)

var (
	usage string = fmt.Sprintf(dedent.Dedent(`
		Usage:
		  dispatch [options]

		Options:
		  -c, --config=FILE      configuration file
		  -h, --help             display this help and exit
		  -o, --output=FILE      redirect output to file
		  -P, --procs=(+)PROCS   number of processes (default %d)
		  -v, --verbose          verbose mode
		      --version          show version

		The number of processes is limited to the number of CPU cores available
		locally by default. In a remote execution context, where the number of
		processes must not rely on the local machine, the sign "+" can be used to
		by-pass this limitation. For example, "dispatch -P +16" will spawn 16
		processes regardless of the number of CPU cores available locally.
 	`)[1:], config.ProcessesDefault)

	out *os.File = os.Stderr
	err error
)

func Dispatch(version string) {
	setEnvirons()
	setupLogging(out, false)

	f := parseFlags()
	k := koanf.New(".")

	// load from the flag set
	if err := config.LoadFlags(k, f); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if k.Bool("version") {
		fmt.Println(version)
		return
	}

	// load from the YAML defined by the config flag
	if err := config.LoadYaml(k, k.String("config")); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// redirect path to a file if specified
	if path := k.String("output"); path != "" {
		out, err = openOutputFile(path)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}

	setupLogging(out, k.Bool("verbose"))

	cfg, err := config.New(k.String("config"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	t, err := cfg.Tasks()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if len(t) == 0 {
		slog.Error("no task to perform")
		os.Exit(1)
	}

	procs := config.ValidateProcs(k.Int("procs"), k.Bool("remote"))
	dispatcher := routines.NewLeader(procs)

	for _, t := range t {
		dispatcher.AddTask(t)
	}

	slog.Info(
		"loading configuration",
		"tasks", len(t),
		"procs", procs,
		"remote", k.Bool("remote"),
		"verbose", k.Bool("verbose"),
	)

	dispatcher.Wait()
	os.Exit(0)
}

func setEnvirons() {
	os.Setenv("PGAPPNAME", "dispatch")
}

func parseFlags() *flag.FlagSet {
	f := config.Flags()
	f.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}

	f.Parse(os.Args[1:])
	return f
}
