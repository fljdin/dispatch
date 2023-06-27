package actions_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks/actions"
	"github.com/stretchr/testify/require"
)

func TestCommandBasicRun(t *testing.T) {
	r := require.New(t)

	cmd := Command{
		Type: "sh",
		Text: "echo test",
	}
	result, _ := cmd.Run()

	r.Equal(OK, result.Status)
	r.Contains(result.Output, "test")
}

func TestCommandValidate(t *testing.T) {
	r := require.New(t)

	cmd := Command{
		Type: "unknown",
		Text: "unknown",
	}
	err := cmd.Validate()

	r.NotNil(err)
	r.Equal("unknown is not supported", err.Error())
}

func TestCommandWithOutput(t *testing.T) {
	r := require.New(t)

	cmd := Command{Text: "echo test"}
	result, _ := cmd.Run()

	r.Equal(OK, result.Status)
	r.Contains(result.Output, "test")
}

func TestCommandWithError(t *testing.T) {
	r := require.New(t)

	cmd := Command{Text: "false"}
	result, _ := cmd.Run()

	r.Equal(KO, result.Status)
}
