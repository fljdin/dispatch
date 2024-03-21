package internal

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func setupLogging() {
	var (
		h     slog.Handler
		level slog.Level = slog.LevelInfo
	)

	var levelStrings = map[string]string{
		// Colors from journalctl.
		"DEBUG": "\x1b[2mDEBUG\x1b[0m ",
		"INFO":  "\x1b[1mINFO\x1b[0m  ",
		"WARN":  "\x1b[1;38;5;185mWARN\x1b[0m  ",
		"ERROR": "\x1b[1;31mERROR\x1b[0m ",
	}

	if argVerbose {
		level = slog.LevelDebug
	}

	w := os.Stderr
	h = tint.NewHandler(w, &tint.Options{
		Level:      level,
		TimeFormat: time.DateTime,
		NoColor:    !isatty.IsTerminal(w.Fd()),

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.LevelKey:
				a.Value = slog.StringValue(levelStrings[a.Value.String()])
			}

			return a
		},
	})

	slog.SetDefault(slog.New(h))
}
