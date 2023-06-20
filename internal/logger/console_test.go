package logger_test

import (
	"testing"
	"time"

	. "github.com/fljdin/dispatch/internal/logger"
	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsoleRender(t *testing.T) {
	console := Console{}
	data, err := console.Parse(Result{
		ID:       1,
		WorkerID: 1,
		SubID:    0,
		Status:   Succeeded,
		Elapsed:  time.Duration(1e8),
	})

	require.Nil(t, err)
	expected := "Worker 1 completed Task 1 (query #0) (success: true, elapsed: 100ms)\n"
	assert.Equal(t, expected, data)
}
