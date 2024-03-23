package logger

import (
	"fmt"
	"log/slog"

	"github.com/fljdin/dispatch/internal/tasks"
)

type Console struct{}

func (c *Console) Render(result tasks.Result) error {

	if !tasks.IsSucceeded(result.Status) {
		slog.Error(result.Code(), result.LoggerArgs()...)
	} else {
		slog.Info(result.Code(), result.LoggerArgs()...)
	}

	slog.Debug(result.Code(), "action", result.Action)

	if len(result.Error) > 0 {
		msg := fmt.Sprintf("%s Error:\n%s", result.Code(), result.Error)
		slog.Error(msg)
	}

	if len(result.Output) > 0 {
		msg := fmt.Sprintf("%s Output:\n%s", result.Code(), result.Output)
		slog.Debug(msg)
	}

	return nil
}
