package dispatcher_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/dispatcher"
	. "github.com/fljdin/dispatch/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestObserverSummaryTrace(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "trace_*.out")

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	observer := &Observer{}
	observer.TraceTo(tempFile.Name())

	observer.Log(TaskResult{
		ID:     1,
		Status: Succeeded,
		Output: "test\n",
	})

	data, err := os.ReadFile(tempFile.Name())
	if assert.Equal(t, nil, err) {
		assert.Contains(t, string(data), "test\n")
	}
}
