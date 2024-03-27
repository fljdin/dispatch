package dispatcher_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/dispatcher"
	. "github.com/fljdin/dispatch/internal/status"
	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/require"
)

func TestDispatcherAddTask(t *testing.T) {
	r := require.New(t)

	dispatcher := New(1)
	dispatcher.AddTask(Task{
		Identifier: NewId(1, 0),
		Action:     Command{Text: "true"},
	})

	r.Equal(Waiting, dispatcher.Evaluate(1))
}

func TestDispatcherDependentTaskNeverExecuted(t *testing.T) {
	r := require.New(t)

	dispatcher := New(1)
	dispatcher.AddTask(Task{
		Identifier: NewId(1, 0),
		Action:     Command{Text: "false"},
	})
	dispatcher.AddTask(Task{
		Identifier: NewId(2, 0),
		Depends:    []int{1},
		Action:     Command{Text: "true"},
	})
	dispatcher.Wait()

	r.Equal(Failed, dispatcher.Evaluate(1))
	r.Equal(Failed, dispatcher.Evaluate(2))
}

func TestDispatcherDependentTaskGetSucceeded(t *testing.T) {
	r := require.New(t)

	dispatcher := New(1)
	dispatcher.AddTask(Task{
		Identifier: NewId(1, 0),
		Action:     Command{Text: "true"},
	})
	dispatcher.AddTask(Task{
		Identifier: NewId(2, 0),
		Depends:    []int{1},
		Action:     Command{Text: "true"},
	})
	dispatcher.Wait()

	r.Equal(Succeeded, dispatcher.Evaluate(1))
	r.Equal(Succeeded, dispatcher.Evaluate(2))
}

func TestDispatcherStatusOfFileTaskMustSummarizeLoadedTaskStatus(t *testing.T) {
	r := require.New(t)

	dispatcher := New(1)
	dispatcher.AddTask(Task{
		Identifier: NewId(1, 0),
		Action:     Command{Text: "false"},
	})
	dispatcher.AddTask(Task{
		Identifier: NewId(1, 1),
		Action:     Command{Text: "true"},
	})
	dispatcher.Wait()

	r.Equal(Failed, dispatcher.Evaluate(1))
}

func TestDispatcherWithOutputLoader(t *testing.T) {
	r := require.New(t)

	dispatcher := New(1)
	dispatcher.AddTask(Task{
		Identifier: NewId(1, 0),
		Action: OutputLoader{
			Text: `echo -n "true\nfalse"`,
			From: "sh",
		},
	})
	dispatcher.Wait()

	r.Equal(Failed, dispatcher.Evaluate(1))
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
		Identifier: NewId(1, 0),
		Action: FileLoader{
			File: tempFile.Name(),
		},
	})
	dispatcher.Wait()

	r.Equal(Failed, dispatcher.Evaluate(1))
}
