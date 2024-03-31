package routines_test

import (
	"os"
	"testing"

	"github.com/fljdin/dispatch/internal/actions"
	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/routines"
	"github.com/fljdin/dispatch/internal/status"
	"github.com/stretchr/testify/require"
)

func TestLeaderAddTask(t *testing.T) {
	r := require.New(t)

	leader := routines.NewLeader(1)
	leader.AddTask(config.Task{
		Identifier: config.NewId(1, 0),
		Action:     actions.Command{Text: "true"},
	})

	r.Equal(status.Waiting, leader.Evaluate(1))
}

func TestLeaderDependentTaskNeverExecuted(t *testing.T) {
	r := require.New(t)

	leader := routines.NewLeader(1)
	leader.AddTask(config.Task{
		Identifier: config.NewId(1, 0),
		Action:     actions.Command{Text: "false"},
	})
	leader.AddTask(config.Task{
		Identifier: config.NewId(2, 0),
		Depends:    []int{1},
		Action:     actions.Command{Text: "true"},
	})
	leader.AddTask(config.Task{
		Identifier: config.NewId(3, 0),
		Depends:    []int{2},
		Action:     actions.Command{Text: "true"},
	})
	leader.Wait()

	r.Equal(status.Failed, leader.Evaluate(1))
	r.Equal(status.Interrupted, leader.Evaluate(2))
	r.Equal(status.Interrupted, leader.Evaluate(2))
}

func TestLeaderDependentTaskGetSucceeded(t *testing.T) {
	r := require.New(t)

	leader := routines.NewLeader(1)
	leader.AddTask(config.Task{
		Identifier: config.NewId(1, 0),
		Action:     actions.Command{Text: "true"},
	})
	leader.AddTask(config.Task{
		Identifier: config.NewId(2, 0),
		Depends:    []int{1},
		Action:     actions.Command{Text: "true"},
	})
	leader.Wait()

	r.Equal(status.Succeeded, leader.Evaluate(1))
	r.Equal(status.Succeeded, leader.Evaluate(2))
}

func TestLeaderStatusOfFileTaskMustSummarizeLoadedTaskStatus(t *testing.T) {
	r := require.New(t)

	leader := routines.NewLeader(1)
	leader.AddTask(config.Task{
		Identifier: config.NewId(1, 0),
		Action:     actions.Command{Text: "false"},
	})
	leader.AddTask(config.Task{
		Identifier: config.NewId(1, 1),
		Action:     actions.Command{Text: "true"},
	})
	leader.Wait()

	r.Equal(status.Failed, leader.Evaluate(1))
}

func TestLeaderWithOutputLoader(t *testing.T) {
	r := require.New(t)

	leader := routines.NewLeader(1)
	leader.AddTask(config.Task{
		Identifier: config.NewId(1, 0),
		Action: actions.OutputLoader{
			From: actions.Shell,
			Text: `echo -n "true\nfalse"`,
		},
	})
	leader.Wait()

	r.Equal(status.Failed, leader.Evaluate(1))
}

func TestLeaderWithFileLoader(t *testing.T) {
	r := require.New(t)

	shFilename := "commands_*.sh"
	shContent := `true\nfalse`
	tempFile, _ := os.CreateTemp("", shFilename)
	tempFile.Write([]byte(shContent))

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	leader := routines.NewLeader(1)
	leader.AddTask(config.Task{
		Identifier: config.NewId(1, 0),
		Action: actions.FileLoader{
			File: tempFile.Name(),
		},
	})
	leader.Wait()

	r.Equal(status.Failed, leader.Evaluate(1))
}
