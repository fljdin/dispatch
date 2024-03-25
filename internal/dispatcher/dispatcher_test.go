package dispatcher_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/dispatcher"
	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/require"
)

func TestDispatcherAddTask(t *testing.T) {
	r := require.New(t)

	dispatcher := New(1)
	dispatcher.AddTask(Task{
		ID:     1,
		Action: Command{Text: "true"},
	})

	r.Equal(Waiting, dispatcher.Status(1))
}

func TestDispatcherDependentTaskNeverExecuted(t *testing.T) {
	r := require.New(t)

	dispatcher := New(1)
	dispatcher.AddTask(Task{
		ID:     1,
		Action: Command{Text: "false"},
	})
	dispatcher.AddTask(Task{
		ID:      2,
		Depends: []int{1},
		Action:  Command{Text: "true"},
	})
	dispatcher.Wait()

	r.Equal(Failed, dispatcher.Status(1))
	r.Equal(Interrupted, dispatcher.Status(2))
}

func TestDispatcherDependentTaskGetSucceeded(t *testing.T) {
	r := require.New(t)

	dispatcher := New(1)
	dispatcher.AddTask(Task{
		ID:     1,
		Action: Command{Text: "true"},
	})
	dispatcher.AddTask(Task{
		ID:      2,
		Depends: []int{1},
		Action:  Command{Text: "true"},
	})
	dispatcher.Wait()

	r.Equal(Succeeded, dispatcher.Status(1))
	r.Equal(Succeeded, dispatcher.Status(2))
}

func TestDispatcherStatusOfFileTaskMustSummarizeLoadedTaskStatus(t *testing.T) {
	r := require.New(t)

	dispatcher := New(1)
	dispatcher.AddTask(Task{
		ID:     1,
		SubID:  0,
		Action: Command{Text: "false"},
	})
	dispatcher.AddTask(Task{
		ID:     1,
		SubID:  1,
		Action: Command{Text: "true"},
	})
	dispatcher.Wait()

	r.Equal(Failed, dispatcher.Status(1))
}

func TestDispatcherWithOutputLoader(t *testing.T) {
	r := require.New(t)

	dispatcher := New(1)
	dispatcher.AddTask(Task{
		ID: 1,
		Action: OutputLoader{
			Text: `echo -n "true\nfalse"`,
			From: "sh",
		},
	})
	dispatcher.Wait()

	r.Equal(Failed, dispatcher.Status(1))
}

func TestDispatcherWithFileLoader(t *testing.T) {
	r := require.New(t)

	shFilename := "commands_*.sh"
	shContent := `true\nfalse`
	tempFile, _ := os.CreateTemp("", shFilename)
	tempFile.Write([]byte(shContent))

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	dispatcher := New(1)
	dispatcher.AddTask(Task{
		ID: 1,
		Action: FileLoader{
			File: tempFile.Name(),
		},
	})
	dispatcher.Wait()

	r.Equal(Failed, dispatcher.Status(1))
}
