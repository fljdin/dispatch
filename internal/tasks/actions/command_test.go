package actions_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandBasicRun(t *testing.T) {
	cmd := Command{
		Type: "sh",
		Text: "echo test",
	}
	result, _ := cmd.Run()

	assert.Equal(t, OK, result.Status)
	assert.Contains(t, result.Output, "test")
}

func TestCommandValidate(t *testing.T) {
	cmd := Command{
		Type: "unknown",
		Text: "unknown",
	}
	err := cmd.Validate()

	require.NotNil(t, err)
	assert.Equal(t, "unknown is not supported", err.Error())
}

func TestCommandWithOutput(t *testing.T) {
	cmd := Command{Text: "echo test"}
	result, _ := cmd.Run()

	assert.Equal(t, OK, result.Status)
	assert.Contains(t, result.Output, "test")
}

func TestCommandWithError(t *testing.T) {
	cmd := Command{Text: "false"}
	result, _ := cmd.Run()

	assert.Equal(t, KO, result.Status)
}
