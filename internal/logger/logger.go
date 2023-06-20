package logger

import (
	"text/template"
	"time"

	"github.com/fljdin/dispatch/internal/task"
)

func newTemplate(name string) *template.Template {
	return template.New(name).Funcs(
		template.FuncMap{
			"isSucceeded": func(status int) bool {
				return status == task.Succeeded
			},
			"roundToMilliseconds": func(duration time.Duration) time.Duration {
				return duration.Round(time.Millisecond)
			},
		},
	)
}

type Logger interface {
	Parse(result task.Result) (string, error)
	Render(result task.Result) error
}
