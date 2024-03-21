package logger

import (
	"github.com/fljdin/dispatch/internal/tasks"
)

type Logger interface {
	Render(result tasks.Result) error
}
