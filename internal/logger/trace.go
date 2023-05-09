package logger

import (
	"os"

	"github.com/fljdin/dispatch/internal/models"
)

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

func (t *Trace) Render(result models.TaskResult) {
	tmpl := newTemplate("trace")
	tmpl, err := tmpl.Parse(TraceTemplate)

	if err != nil {
		panic(err)
	}

	tmpl.Execute(t.file, result)
}
