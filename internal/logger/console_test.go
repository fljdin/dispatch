package logger_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	. "github.com/fljdin/dispatch/internal/logger"
	. "github.com/fljdin/dispatch/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestConsoleRender(t *testing.T) {
	r, w, _ := os.Pipe()
	tmp := os.Stdout
	os.Stdout = w

	defer r.Close()
	defer func() { os.Stdout = tmp }()

	console := &Console{}
	console.Render(TaskResult{
		ID:       1,
		WorkerID: 1,
		QueryID:  0,
		Status:   Succeeded,
		Elapsed:  time.Duration(1e8),
	})

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	expected := "Worker 1 completed Task 1 (query #0) (success: true, elapsed: 100ms)\n"
	actual := buf.String()
	assert.Equal(t, expected, actual)
}
