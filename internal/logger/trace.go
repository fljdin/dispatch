package logger

import (
	"bytes"
	"os"
	"text/template"
	"time"

	"github.com/fljdin/dispatch/internal/tasks"
)

func newTemplate(name string) *template.Template {
	return template.New(name).Funcs(
		template.FuncMap{
			"isSucceeded": func(status int) bool {
				return tasks.IsSucceeded(status)
			},
			"roundToMilliseconds": func(duration time.Duration) time.Duration {
				return duration.Round(time.Millisecond)
			},
		},
	)
}

const TraceTemplate string = `===== Task {{.ID}} (command #{{.SubID}}) (success: {{if isSucceeded .Status}}true{{else}}false{{end}}, elapsed: {{roundToMilliseconds .Elapsed}}) =====
Started at: {{.StartTime}}
Ended at:   {{.EndTime}}
Error: {{.Error}}
Output:
{{.Output}}
`

type Trace struct {
	Filename string
	file     *os.File
}

func (t *Trace) Open() error {
	var err error
	const flag int = os.O_APPEND | os.O_TRUNC | os.O_CREATE | os.O_WRONLY

	t.file, err = os.OpenFile(t.Filename, flag, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (t *Trace) Parse(result tasks.Result) (string, error) {
	tmpl := newTemplate("trace")
	tmpl, err := tmpl.Parse(TraceTemplate)

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	tmpl.Execute(&buf, result)

	return buf.String(), nil
}

func (t *Trace) Render(result tasks.Result) error {
	data, err := t.Parse(result)
	if err != nil {
		return nil
	}

	_, err = t.file.Write([]byte(data))
	return err
}
