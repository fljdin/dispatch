package tasks

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type Task struct {
	ID      int
	SubID   int
	Name    string
	Action  Action
	Depends []int
	Status  int
}

func (t Task) Validate() error {
	if t.ID == 0 {
		return fmt.Errorf("id is required")
	}

	if t.Action == nil {
		return fmt.Errorf("action is required")
	}

	if err := t.Action.Validate(); err != nil {
		return err
	}

	return nil
}

func (t Task) ValidateDependencies(identifiers []int) error {
	exists := true

	for _, d := range t.Depends {
		exists = exists && slices.Contains(identifiers, d)
	}

	if !exists {
		return fmt.Errorf("task %d depends on unknown task %d", t.ID, t.Depends)
	}

	return nil
}
