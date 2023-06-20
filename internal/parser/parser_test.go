package parser_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserWithSqlContent(t *testing.T) {
	sqlContent := "SELECT 1; SELECT 2; SELECT 3;"

	parser, _ := NewBuilder("psql").
		WithContent(sqlContent).
		Build()

	queries := parser.Parse()
	assert.Equal(t, 3, len(queries))
}

func TestParserHandleStrings(t *testing.T) {
	sqlContent := []string{
		`SELECT ';"';`,
		`SELECT 1 ";'";`,
		"SELECT $$;$$;",
		"SELECT $tag$;$tag$;",
		"SELECT $tag$$tag;$tag$;",
	}

	for i, q := range sqlContent {
		parser, _ := NewBuilder("psql").
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
		parser, _ := NewBuilder("psql").
			WithContent(q).
			Build()

		queries := parser.Parse()
		assert.Equal(t, sqlContent[i], queries[0])
	}
}

func TestParserHandleCommentsAndStringsMix(t *testing.T) {
	sqlContent := []string{
		`SELECT /*'*/ 1"';";`,
		`SELECT $$/*$$ AS "$$;";`,
		"SELECT 1 -- $tag$ ;\n +2;",
		"SELECT /* $tag$ */$tag$;$tag$;",
	}

	for i, q := range sqlContent {
		parser, _ := NewBuilder("psql").
			WithContent(q).
			Build()

		queries := parser.Parse()
		assert.Equal(t, sqlContent[i], queries[0])
	}
}

func TestParserHandleTransactionBlock(t *testing.T) {
	sqlContent := []string{
		"BEGIN; SELECT 1; END;",
		"BEGIN; SELECT 1; COMMIT;",
		"BEGIN; SELECT 1; ROLLBACK;",
		"begin; SELECT 'END'; end;",
	}

	for i, q := range sqlContent {
		parser, _ := NewBuilder("psql").
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

	parser, _ := NewBuilder("psql").
		FromFile(tempFile.Name()).
		Build()

	queries := parser.Parse()
	assert.Equal(t, 1, len(queries))
}

func TestParserGCommand(t *testing.T) {
	sqlContent := `SELECT 1\g
SELECT 2\g result.txt
SELECT 3\g (format=unaligned tuples_only)
`
	parser, _ := NewBuilder("psql").
		WithContent(sqlContent).
		Build()

	queries := parser.Parse()
	require.Equal(t, 3, len(queries))

	assert.Equal(t, "SELECT 1\\g\n", queries[0])
	assert.Equal(t, "SELECT 2\\g result.txt\n", queries[1])
	assert.Equal(t, "SELECT 3\\g (format=unaligned tuples_only)\n", queries[2])
}

func TestParserCrosstabviewCommand(t *testing.T) {
	sqlContent := `SELECT 1, 1, 1 \crosstabview
SELECT 2, 2, 2 \crosstabview
SELECT 3, 3, 3 \crosstabview
`

	parser, _ := NewBuilder("psql").
		WithContent(sqlContent).
		Build()

	queries := parser.Parse()
	assert.Equal(t, 3, len(queries))
}

func TestParserUnsupportedCommand(t *testing.T) {
	sqlContent := "SELECT 1\\unsupported\nSELECT 1\\\nSELECT 1;"

	parser, _ := NewBuilder("psql").
		WithContent(sqlContent).
		Build()

	queries := parser.Parse()
	assert.Equal(t, 1, len(queries))
}
