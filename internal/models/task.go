package models

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"time"

	"golang.org/x/exp/slices"
)

var TaskTypes = []string{"sh", "psql"}

type Task struct {
	ID         int    `yaml:"id"`
	Type       string `yaml:"type,omitempty"`
	Name       string `yaml:"name,omitempty"`
	Command    string `yaml:"command"`
	File       string `yaml:"file"`
	URI        string `yaml:"uri,omitempty"`
	Connection string `yaml:"connection,omitempty"`
	Depends    []int  `yaml:"depends_on,omitempty"`
	Output     string `yaml:"output,omitempty"`
	QueryID    int
}

func (t Task) VerifyRequired() error {
	if t.ID == 0 {
		return fmt.Errorf("id is required")
	}

	if t.Command == "" && t.File == "" {
		return fmt.Errorf("command is required")
	}

	return nil
}

func (t Task) VerifyType() error {
	if t.Type == "" {
		return nil
	}
	if !slices.Contains(TaskTypes, t.Type) {
		return fmt.Errorf("%s is an invalid task type", t.Type)
	}
	return nil
}

func (t Task) VerifyDependencies(identifiers []int) error {
	verified := true

	for _, d := range t.Depends {
		verified = verified && slices.Contains(identifiers, d)
	}

	if !verified {
		return fmt.Errorf("task %d depends on unknown task %d", t.ID, t.Depends)
	}

	return nil
}

func (t Task) VerifyOutput() error {
	if t.Output == "" {
		return nil
	}

	_, err := template.New("ouptut").Parse(t.Output)
	if err != nil {
		return err
	}

	return nil
}

func (t Task) GetOutput() string {
	var buf bytes.Buffer

	tmpl, _ := template.New("ouptut").Parse(t.Output)
	tmpl.Execute(&buf, t)

	return buf.String()
}

func (t Task) writeOutput(output []byte) error {
	if len(t.Output) == 0 {
		return nil
	}
	return os.WriteFile(t.GetOutput(), output, 0644)
}

func (t Task) Run() TaskResult {
	var cmd *exec.Cmd

	startTime := time.Now()

	switch t.Type {
	case "psql":
		cmd = exec.Command("psql", "-d", t.URI, "-c", t.Command)
	default:
		cmd = exec.Command("sh", "-c", t.Command)
	}

	output, err := cmd.CombinedOutput()
	endTime := time.Now()

	tr := TaskResult{
		ID:        t.ID,
		QueryID:   t.QueryID,
		StartTime: startTime,
		EndTime:   endTime,
		Elapsed:   endTime.Sub(startTime),
		Status:    Succeeded,
		Output:    string(output),
	}

	if err != nil {
		tr.Status = Failed
		tr.Error = err.Error()
	}

	if err = t.writeOutput(output); err != nil {
		if len(tr.Error) > 0 {
			tr.Error += "\n" + err.Error()
		} else {
			tr.Error = err.Error()
		}
	}

	return tr
}
