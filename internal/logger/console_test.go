package logger_test

import (
	"testing"
	"time"

	. "github.com/fljdin/dispatch/internal/logger"
	. "github.com/fljdin/dispatch/internal/task"
	"github.com/stretchr/testify/assert"
)

func TestConsoleRender(t *testing.T) {
	console := Console{}
	data, _ := console.Parse(Result{
		ID:       1,
		WorkerID: 1,
		QueryID:  0,
		Status:   Succeeded,
		Elapsed:  time.Duration(1e8),
	})

	expected := "Worker 1 completed Task 1 (query #0) (success: true, elapsed: 100ms)\n"
	assert.Equal(t, expected, data)
}
