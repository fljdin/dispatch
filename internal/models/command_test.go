package models_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandBasicRun(t *testing.T) {
	cmd := Command{
		Type: "sh",
		Text: "echo test",
	}
	result := cmd.Run()

	assert.Equal(t, Succeeded, result.Status)
	assert.Contains(t, result.Output, "test")
}

func TestCommandVerifyType(t *testing.T) {
	cmd := Command{
		Type: "unknown",
		Text: "unknown",
	}
	err := cmd.VerifyType()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "is not supported")
}

func TestCommandWithOutput(t *testing.T) {
	cmd := Command{Text: "echo test"}
	result := cmd.Run()

	assert.Equal(t, Succeeded, result.Status)
	assert.Contains(t, result.Output, "test")
}

func TestCommandWithError(t *testing.T) {
	cmd := Command{Text: "false"}
	result := cmd.Run()

	assert.Equal(t, Failed, result.Status)
}
