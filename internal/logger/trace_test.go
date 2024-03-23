package logger_test

import (
	"testing"
	"time"

	. "github.com/fljdin/dispatch/internal/logger"
	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/require"
)

func TestTraceRender(t *testing.T) {
	r := require.New(t)

	trace := Trace{Filename: "dummy.txt"}
	data, err := trace.Parse(Result{
		ID:        1,
		ProcID:    1,
		SubID:     0,
		Status:    Succeeded,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Elapsed:   time.Duration(1e8),
		Output:    "test\n",
	})

	r.Nil(err)
	r.Contains(data, "Output:\ntest\n")
}
