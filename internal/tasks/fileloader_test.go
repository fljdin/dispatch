package tasks_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/status"
	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/require"
)

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
		Type: PgSQL,
	}
	result, commands := cmd.Run()

	r.Equal(Succeeded, result.Status)
	r.Equal(Command{Text: "SELECT 1;", Type: PgSQL}, commands[0])
	r.Equal(Command{Text: "SELECT 2;", Type: PgSQL}, commands[1])
}
