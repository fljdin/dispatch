package logger_test

import (
	"testing"
	"time"

	. "github.com/fljdin/dispatch/internal/logger"
	. "github.com/fljdin/dispatch/internal/task"
	"github.com/stretchr/testify/assert"
)

func TestTraceRender(t *testing.T) {
	trace := Trace{Filename: "dummy.txt"}
	data, _ := trace.Parse(TaskResult{
		ID:        1,
		WorkerID:  1,
		QueryID:   0,
		Status:    Succeeded,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Elapsed:   time.Duration(1e8),
		Output:    "test\n",
	})

	assert.Contains(t, data, "Output:\ntest\n")
}
