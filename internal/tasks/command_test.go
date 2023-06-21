package tasks_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks"
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

func TestCommandVerifyRequired(t *testing.T) {
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
	result := cmd.Run()

	assert.Equal(t, Succeeded, result.Status)
	assert.Contains(t, result.Output, "test")
}

func TestCommandWithError(t *testing.T) {
	cmd := Command{Text: "false"}
	result := cmd.Run()

	assert.Equal(t, Failed, result.Status)
}

func TestCommandWithInvalidGeneratorType(t *testing.T) {
	cmd := Command{
		Text: "unknown",
		From: "unknown",
	}
	result, _ := cmd.Generate()

	assert.Equal(t, Failed, result.Status)
	assert.Equal(t, "unknown is not supported", result.Error)
}

func TestCommandWithFailedGenerator(t *testing.T) {
	cmd := Command{
		Text: `echo true ; false`,
		From: "sh",
	}
	result, _ := cmd.Generate()

	assert.Equal(t, Failed, result.Status)
}

func TestCommandGeneratorFromText(t *testing.T) {
	cmd := Command{
		Text: `echo -n "true\nfalse"`,
		From: "sh",
	}
	result, commands := cmd.Generate()

	require.Equal(t, Succeeded, result.Status)
	assert.Equal(t, Command{Text: "true", Type: "sh"}, commands[0])
	assert.Equal(t, Command{Text: "false", Type: "sh"}, commands[1])
}

func TestCommandGeneratorFromFile(t *testing.T) {
	sqlFilename := "queries_*.sql"
	sqlContent := "SELECT 1;SELECT 2;"
	tempFile, _ := os.CreateTemp("", sqlFilename)
	tempFile.Write([]byte(sqlContent))

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	cmd := Command{
		File: tempFile.Name(),
		From: "psql",
	}
	result, commands := cmd.Generate()

	require.Equal(t, Succeeded, result.Status)
	assert.Equal(t, Command{Text: "SELECT 1;", Type: "psql"}, commands[0])
	assert.Equal(t, Command{Text: "SELECT 2;", Type: "psql"}, commands[1])
}
