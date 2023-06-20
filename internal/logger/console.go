package logger

import (
	"bytes"
	"log"

	"github.com/fljdin/dispatch/internal/task"
)

const ConsoleTemplate string = `Worker {{.WorkerID}} completed Task {{.ID}} (query #{{.QueryID}}) (success: {{if isSucceeded .Status}}true{{else}}false{{end}}, elapsed: {{roundToMilliseconds .Elapsed}})
`

type Console struct{}

func (c *Console) Parse(result task.TaskResult) (string, error) {
	tmpl := newTemplate("console")
	tmpl, err := tmpl.Parse(ConsoleTemplate)

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	tmpl.Execute(&buf, result)

	return buf.String(), nil
}

func (c *Console) Render(result task.TaskResult) error {
	data, err := c.Parse(result)

	if err != nil {
		return err
	}

	log.Print(data)
	return nil
}
