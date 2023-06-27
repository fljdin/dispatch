package actions_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks/actions"
	"github.com/stretchr/testify/require"
)

func TestOutputLoaderWithInvalidType(t *testing.T) {
	r := require.New(t)

	cmd := OutputLoader{
		Text: "unknown",
		From: "unknown",
	}
	result, _ := cmd.Run()

	r.Equal(KO, result.Status)
	r.Equal("unknown is not supported", result.Error)
}

func TestOutputLoaderWithFailedCommand(t *testing.T) {
	r := require.New(t)

	cmd := OutputLoader{
		Text: `echo true ; false`,
		From: "sh",
	}
	result, _ := cmd.Run()

	r.Equal(KO, result.Status)
}

func TestOutputLoaderRun(t *testing.T) {
	r := require.New(t)

	cmd := OutputLoader{
		Text: "echo true; echo false",
		From: "sh",
	}
	result, commands := cmd.Run()

	r.Equal(OK, result.Status)
	r.Equal("true\nfalse\n", result.Output)
	r.Equal(Command{Text: "true"}, commands[0])
	r.Equal(Command{Text: "false"}, commands[1])
}
