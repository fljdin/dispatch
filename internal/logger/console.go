package logger

import (
	"log/slog"

	"github.com/fljdin/dispatch/internal/tasks"
)

type Console struct{}

func (c *Console) Render(result tasks.Result) error {
	if !tasks.IsSucceeded(result.Status) {
		slog.Error(result.Error, result.LoggerArgs()...)
	} else {
		slog.Info("task completed", result.LoggerArgs()...)
	}

	return nil
}
