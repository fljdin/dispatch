package dispatcher_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/dispatcher"
	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/assert"
)

func TestDispatcherAddTask(t *testing.T) {
	dispatcher, _ := NewBuilder().Build()

	dispatcher.AddTask(Task{
		ID:      1,
		Command: Command{Text: "true"},
	})

	assert.Equal(t, Waiting, dispatcher.Status(1))
}

func TestDispatcherDependentTaskNeverExecuted(t *testing.T) {
	dispatcher, _ := NewBuilder().Build()

	dispatcher.AddTask(Task{
		ID:      1,
		Command: Command{Text: "false"},
	})
	dispatcher.AddTask(Task{
		ID:      2,
		Depends: []int{1},
		Command: Command{Text: "true"},
	})
	dispatcher.Wait()

	assert.Equal(t, Failed, dispatcher.Status(1))
	assert.Equal(t, Interrupted, dispatcher.Status(2))
}

func TestDispatcherDependentTaskGetSucceeded(t *testing.T) {
	dispatcher, _ := NewBuilder().Build()

	dispatcher.AddTask(Task{
		ID:      1,
		Command: Command{Text: "true"},
	})
	dispatcher.AddTask(Task{
		ID:      2,
		Depends: []int{1},
		Command: Command{Text: "true"},
	})
	dispatcher.Wait()

	assert.Equal(t, Succeeded, dispatcher.Status(1))
	assert.Equal(t, Succeeded, dispatcher.Status(2))
}

func TestDispatcherStatusOfFileTaskMustSummarizeLoadedTaskStatus(t *testing.T) {
	dispatcher, _ := NewBuilder().Build()

	dispatcher.AddTask(Task{
		ID:      1,
		SubID:   0,
		Command: Command{Text: "false"},
	})
	dispatcher.AddTask(Task{
		ID:      1,
		SubID:   1,
		Command: Command{Text: "true"},
	})
	dispatcher.Wait()

	assert.Equal(t, Failed, dispatcher.Status(1))
}

func TestDispatcherGeneratedTaskFromCommand(t *testing.T) {
	dispatcher, _ := NewBuilder().Build()

	dispatcher.AddTask(Task{
		ID: 1,
		Command: Command{
			Text: `echo -n "true\nfalse"`,
			From: "sh",
		},
	})
	dispatcher.Wait()

	assert.Equal(t, Failed, dispatcher.Status(1))
}

func TestDispatcherGenerateTaskFromFile(t *testing.T) {
	shFilename := "commands_*.sh"
	shContent := `true\nfalse`
	tempFile, _ := os.CreateTemp("", shFilename)
	tempFile.Write([]byte(shContent))

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	dispatcher, _ := NewBuilder().Build()
	dispatcher.AddTask(Task{
		ID: 1,
		Command: Command{
			File: tempFile.Name(),
		},
	})

	dispatcher.Wait()
	assert.Equal(t, Failed, dispatcher.Status(1))
}
