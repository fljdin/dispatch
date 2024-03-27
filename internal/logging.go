package internal

import (
	"log/slog"
	"os"
	"time"

	"github.com/fljdin/dispatch/internal/helper"
	"github.com/lmittmann/tint"
)

func openOutputFile(filename string) (*os.File, error) {
	const flag int = os.O_APPEND | os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	return os.OpenFile(filename, flag, 0644)
}

func setupLogging(out *os.File, verbose bool) {
	var (
		h     slog.Handler
		level slog.Level = slog.LevelInfo
	)

	var levelStrings = map[string]string{
		"DEBUG": "DEBUG ",
		"INFO":  "INFO  ",
		"WARN":  "WARN  ",
		"ERROR": "ERROR ",
	}

	if verbose {
		level = slog.LevelDebug
	}

	h = tint.NewHandler(out, &tint.Options{
		Level:      level,
		TimeFormat: time.DateTime,
		NoColor:    true,

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.Attr{
					Key:   "time",
					Value: slog.TimeValue(helper.Now()),
				}
			case slog.LevelKey:
				a.Value = slog.StringValue(levelStrings[a.Value.String()])
			}

			return a
		},
	})

	slog.SetDefault(slog.New(h))
}
