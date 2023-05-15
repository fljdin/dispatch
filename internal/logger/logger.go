package logger

import (
	"text/template"
	"time"

	"github.com/fljdin/dispatch/internal/models"
)

func newTemplate(name string) *template.Template {
	return template.New(name).Funcs(
		template.FuncMap{
			"isSucceeded": func(status int) bool {
				return status == models.Succeeded
			},
			"roundToMilliseconds": func(duration time.Duration) time.Duration {
				return duration.Round(time.Millisecond)
			},
		},
	)
}

type Logger interface {
	Parse(result models.TaskResult) (string, error)
	Render(result models.TaskResult) error
}
