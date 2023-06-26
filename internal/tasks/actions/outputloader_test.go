package actions_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutputLoaderWithInvalidType(t *testing.T) {
	cmd := OutputLoader{
		Text: "unknown",
		From: "unknown",
	}
	result, _ := cmd.Run()

	assert.Equal(t, KO, result.Status)
	assert.Equal(t, "unknown is not supported", result.Error)
}

func TestOutputLoaderWithFailedCommand(t *testing.T) {
	cmd := OutputLoader{
		Text: `echo true ; false`,
		From: "sh",
	}
	result, _ := cmd.Run()

	assert.Equal(t, KO, result.Status)
}

func TestOutputLoaderRun(t *testing.T) {
	cmd := OutputLoader{
		Text: "echo true; echo false",
		From: "sh",
	}
	result, commands := cmd.Run()

	require.Equal(t, OK, result.Status)
	assert.Equal(t, "true\nfalse\n", result.Output)
	assert.Equal(t, Command{Text: "true"}, commands[0])
	assert.Equal(t, Command{Text: "false"}, commands[1])
}
