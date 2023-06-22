package cmd_test

import (
	"os"
	"testing"

	"github.com/fljdin/dispatch/internal/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunWithConfig(t *testing.T) {
	cmd.RootCmd.SetArgs([]string{"run", "--config", "config/TestRunWithConfig.yaml"})
	err := cmd.RootCmd.Execute()
	require.Nil(t, err)

	output, err := os.ReadFile("TestRunWithConfig.log")
	defer os.Remove("TestRunWithConfig.log")
	require.Nil(t, err)

	expected, err := os.ReadFile("expected/TestRunWithConfig.log")
	require.Nil(t, err)

	assert.Equal(t, string(expected), string(output))
}

func TestWorkerForwardURIToGeneratedTasks(t *testing.T) {

}
