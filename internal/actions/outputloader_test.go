package actions_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/actions"
	. "github.com/fljdin/dispatch/internal/status"
	"github.com/stretchr/testify/require"
)

func TestOutputLoaderWithInvalidType(t *testing.T) {
	r := require.New(t)

	cmd := OutputLoader{
		Text: "unknown",
		From: "unknown",
	}
	result, _ := cmd.Run()

	r.Equal(Failed, result.Status)
	r.Equal("unknown is not supported", result.Error)
}

func TestOutputLoaderWithFailedCommand(t *testing.T) {
	r := require.New(t)

	cmd := OutputLoader{
		From: Shell,
		Text: `echo true ; false`,
	}
	result, _ := cmd.Run()

	r.Equal(Failed, result.Status)
}

func TestOutputLoaderRun(t *testing.T) {
	r := require.New(t)

	cmd := OutputLoader{
		From: Shell,
		Text: "echo true; echo false",
	}
	result, commands := cmd.Run()

	r.Equal(Succeeded, result.Status)
	r.Equal("true\nfalse\n", result.Output)
	r.Equal(Command{Text: "true"}, commands[0])
	r.Equal(Command{Text: "false"}, commands[1])
}
