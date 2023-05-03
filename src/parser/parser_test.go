package parser_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/src/parser"
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

func TestParserHandleStrings(t *testing.T) {
	sqlContent := []string{
		`SELECT ';"';`,
		`SELECT 1 ";'";`,
		`SELECT $$;$$;`,
		`SELECT $tag$;$tag$;`,
		`SELECT $tag$$tag;$tag$;`,
	}

	for i, q := range sqlContent {
		parser, _ := NewParserBuilder("psql").
			WithContent(q).
			Build()

		queries := parser.Parse()
		assert.Equal(t, sqlContent[i], queries[0])
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
		assert.Equal(t, sqlContent[i], queries[0])
	}
}

func TestParserHandleTransactionBloc(t *testing.T) {
	sqlContent := []string{
		"BEGIN; SELECT 1; END;",
		"BEGIN; SELECT 1; COMMIT;",
		"BEGIN; SELECT 1; ROLLBACK;",
		"begin; SELECT 'END'; end;",
	}

	for i, q := range sqlContent {
		parser, _ := NewParserBuilder("psql").
			WithContent(q).
			Build()

		queries := parser.Parse()
		assert.Equal(t, sqlContent[i], queries[0])
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
