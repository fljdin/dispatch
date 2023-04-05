package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserWithSqlContent(t *testing.T) {
	sqlContent := "SELECT 1; SELECT 2; SELECT 3;"

	parser, _ := NewParserBuilder("psql").
		WithContent(sqlContent).
		Build()

	queries := parser.Parse()

	assert.Equal(t, len(queries), 3)
}

func TestParserFromSqlFile(t *testing.T) {
	sqlFilename := "queries_*.sql"
	sqlContent := "SELECT 1;"
	tempFile, _ := os.CreateTemp("", sqlFilename)

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	tempFile.Write([]byte(sqlContent))

	parser, _ := NewParserBuilder("psql").
		FromFile(tempFile.Name()).
		Build()
	queries := parser.Parse()

	assert.Equal(t, len(queries), 1)
}

func TestParserFromInvalidParseType(t *testing.T) {
	_, err := NewParserBuilder("unknown").
		Build()

	if assert.NotEqual(t, err, nil) {
		assert.Contains(t, err.Error(), "invalid type for parsing file")
	}
}
