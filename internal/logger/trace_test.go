package logger_test

import (
	"os"
	"testing"
	"time"

	. "github.com/fljdin/dispatch/internal/logger"
	. "github.com/fljdin/dispatch/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTraceRender(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "trace_*.out")

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	trace := &Trace{Filename: tempFile.Name()}
	trace.Open()
	trace.Render(TaskResult{
		ID:        1,
		WorkerID:  1,
		QueryID:   0,
		Status:    Succeeded,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Elapsed:   time.Duration(1e8),
		Output:    "test\n",
	})

	data, err := os.ReadFile(tempFile.Name())
	if assert.Equal(t, nil, err) {
		assert.Contains(t, string(data), "Output:\ntest\n")
	}
}
