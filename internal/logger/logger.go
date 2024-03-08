package logger

import (
	"text/template"
	"time"

	"github.com/fljdin/dispatch/internal/tasks"
	"golang.org/x/exp/slices"
)

func newTemplate(name string) *template.Template {
	return template.New(name).Funcs(
		template.FuncMap{
			"isSucceeded": func(status int) bool {
				s := []int{tasks.Ready, tasks.Succeeded}
				return slices.Contains(s, status)
			},
			"roundToMilliseconds": func(duration time.Duration) time.Duration {
				return duration.Round(time.Millisecond)
			},
		},
	)
}

type Logger interface {
	Parse(result tasks.Result) (string, error)
	Render(result tasks.Result) error
}
