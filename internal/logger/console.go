package logger

import (
	"os"

	"github.com/fljdin/dispatch/internal/models"
)

const ConsoleTemplate string = `Worker {{.WorkerID}} completed Task {{.ID}} (query #{{.QueryID}}) (success: {{if isSucceeded .Status}}true{{else}}false{{end}}, elapsed: {{roundToMilliseconds .Elapsed}})
`

type Console struct{}

func (c *Console) Render(result models.TaskResult) {
	tmpl := newTemplate("console")
	tmpl, err := tmpl.Parse(ConsoleTemplate)

	if err != nil {
		panic(err)
	}

	tmpl.Execute(os.Stdout, result)
}
