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

	assert.Equal(t, 3, len(queries))
}

func TestParserHandleLiterals(t *testing.T) {
	sqlContent := []string{
		`SELECT ';"';`,
		`SELECT 1 ";'";`,
		`SELECT $$;$$;`,
	}

	for i, q := range sqlContent {
		parser, _ := NewParserBuilder("psql").
			WithContent(q).
			Build()

		queries := parser.Parse()
		assert.Equal(t, queries[0], sqlContent[i])
	}
}

func TestParserHandleComments(t *testing.T) {
	sqlContent := []string{
		"SELECT 1 /* comment ; */ + 2;",
		"SELECT 1 /* comment ;\n comment ; */ + 2;",
		"SELECT 1 -- comment ;\n + 2;",
		"SELECT 1 -- /* comment ;\n +2;",
	}

	for i, q := range sqlContent {
		parser, _ := NewParserBuilder("psql").
			WithContent(q).
			Build()

		queries := parser.Parse()
		assert.Equal(t, queries[0], sqlContent[i])
	}
}

func TestParserHandleTransactionBloc(t *testing.T) {
	sqlContent := []string{
		"BEGIN; SELECT 1; END;",
		"BEGIN; SELECT 1; COMMIT;",
		"BEGIN; SELECT 1; ROLLBACK;",
		"BEGIN; SELECT 'END'; END;",
	}

	for i, q := range sqlContent {
		parser, _ := NewParserBuilder("psql").
			WithContent(q).
			Build()

		queries := parser.Parse()
		assert.Equal(t, queries[0], sqlContent[i])
	}
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

	assert.Equal(t, 1, len(queries))
}

func TestParserFromInvalidParseType(t *testing.T) {
	_, err := NewParserBuilder("unknown").
		Build()

	if assert.NotEqual(t, nil, err) {
		assert.Contains(t, err.Error(), "invalid type for parsing file")
	}
}
