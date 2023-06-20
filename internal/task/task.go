package task

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type Task struct {
	ID      int
	Name    string
	Command Command
	Depends []int
	QueryID int
	Status  int
}

func (t Task) VerifyRequired() error {
	if t.ID == 0 {
		return fmt.Errorf("id is required")
	}

	if t.Command.Text == "" && t.Command.File == "" {
		return fmt.Errorf("command is required")
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
