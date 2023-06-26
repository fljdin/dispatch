package logger_test

import (
	"testing"
	"time"

	. "github.com/fljdin/dispatch/internal/logger"
	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTraceRender(t *testing.T) {
	trace := Trace{Filename: "dummy.txt"}
	data, err := trace.Parse(Result{
		ID:        1,
		WorkerID:  1,
		SubID:     0,
		Status:    Succeeded,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Elapsed:   time.Duration(1e8),
		Output:    "test\n",
	})

	require.Nil(t, err)
	assert.Contains(t, data, "Output:\ntest\n")
}
