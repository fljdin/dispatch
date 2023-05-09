package logger

import (
	"text/template"
	"time"

	"github.com/fljdin/dispatch/internal/models"
)

const ConsoleTemplate string = `Worker {{.WorkerID}} completed Task {{.ID}} (query #{{.QueryID}}) (success: {{if isSucceeded .Status}}true{{else}}false{{end}}, elapsed: {{roundToMilliseconds .Elapsed}})
`

const TraceTemplate string = `===== Task {{.ID}} (query #{{.QueryID}}) (success: {{if isSucceeded .Status}}true{{else}}false{{end}}, elapsed: {{roundToMilliseconds .Elapsed}}) =====
Started at: {{.StartTime}}
Ended at:   {{.EndTime}}
Error: {{.Error}}
Output:
{{.Output}}
`

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
	Render(result models.TaskResult)
}
