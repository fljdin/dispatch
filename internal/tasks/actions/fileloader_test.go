package actions_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileLoaderValidate(t *testing.T) {
	cmd := FileLoader{
		File: "unknown.txt",
	}
	err := cmd.Validate()

	require.NotNil(t, err)
	assert.Equal(t, "type is required with a file", err.Error())
}

func TestFileLoaderRun(t *testing.T) {
	sqlFilename := "queries_*.sql"
	sqlContent := "SELECT 1;SELECT 2;"
	tempFile, _ := os.CreateTemp("", sqlFilename)
	tempFile.Write([]byte(sqlContent))

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	cmd := FileLoader{
		File: tempFile.Name(),
		Type: "psql",
	}
	result, commands := cmd.Run()

	require.Equal(t, OK, result.Status)
	assert.Equal(t, Command{Text: "SELECT 1;", Type: "psql"}, commands[0])
	assert.Equal(t, Command{Text: "SELECT 2;", Type: "psql"}, commands[1])
}
