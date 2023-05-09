package logger

import (
	"os"

	"github.com/fljdin/dispatch/internal/models"
)

type Console struct{}

func (c *Console) Render(result models.TaskResult) {
	tmpl := newTemplate("console")
	tmpl, err := tmpl.Parse(ConsoleTemplate)

	if err != nil {
		panic(err)
	}

	tmpl.Execute(os.Stdout, result)
}
