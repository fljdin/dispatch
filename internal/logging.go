package internal

import (
	"log/slog"
	"os"
	"time"

	"github.com/fljdin/dispatch/internal/tasks"
	"github.com/lmittmann/tint"
)

func openOutputFile(filename string) (*os.File, error) {
	const flag int = os.O_APPEND | os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	return os.OpenFile(filename, flag, 0644)
}

func setupLogging(w *os.File) {
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

	if argVerbose {
		level = slog.LevelDebug
	}

	h = tint.NewHandler(w, &tint.Options{
		Level:      level,
		TimeFormat: time.DateTime,
		NoColor:    true,

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.Attr{
					Key:   "time",
					Value: slog.TimeValue(tasks.Time()),
				}
			case slog.LevelKey:
				a.Value = slog.StringValue(levelStrings[a.Value.String()])
			}

			return a
		},
	})

	slog.SetDefault(slog.New(h))
}
