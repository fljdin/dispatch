package tasks

import (
	"fmt"

	"github.com/fljdin/dispatch/internal/status"
	"golang.org/x/exp/slices"
)

type TaskIdentifier struct {
	ID    int
	SubID int
}

func NewId(id, subid int) TaskIdentifier {
	return TaskIdentifier{
		ID:    id,
		SubID: subid,
	}
}

func (t TaskIdentifier) IsZero() bool {
	return t.ID == 0 && t.SubID == 0
}

type Task struct {
	Identifier TaskIdentifier
	Name       string
	Action     Action
	Depends    []int
	Status     status.Status
}

func (t Task) String() string {
	return fmt.Sprintf("%d:%d", t.Identifier.ID, t.Identifier.SubID)
}

func (t Task) Validate() error {
	if t.Identifier.IsZero() {
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
		return fmt.Errorf("task %d depends on unknown task %d", t.Identifier.ID, t.Depends)
	}

	return nil
}
