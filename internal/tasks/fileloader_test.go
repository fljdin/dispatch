package tasks_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/require"
)

func TestFileLoaderValidate(t *testing.T) {
	r := require.New(t)

	cmd := FileLoader{
		File: "unknown.txt",
	}
	err := cmd.Validate()

	r.NotNil(err)
	r.Equal("type is required with a file", err.Error())
}

func TestFileLoaderRun(t *testing.T) {
	r := require.New(t)

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

	r.Equal(Ready, result.Status)
	r.Equal(Command{Text: "SELECT 1;", Type: "psql"}, commands[0])
	r.Equal(Command{Text: "SELECT 2;", Type: "psql"}, commands[1])
}
