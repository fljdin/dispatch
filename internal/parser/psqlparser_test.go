package parser_test

import (
	"os"
	"testing"

	. "github.com/fljdin/dispatch/internal/parser"
	"github.com/stretchr/testify/require"
)

func TestPsqlParserWithContent(t *testing.T) {
	r := require.New(t)

	sqlContent := "SELECT 1; SELECT 2; SELECT 3;"
	parser, _ := NewBuilder("psql").
		WithContent(sqlContent).
		Build()
	queries := parser.Parse()

	r.Equal(3, len(queries))
}

func TestPsqlParserHandleStrings(t *testing.T) {
	r := require.New(t)

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

		r.Equal(sqlContent[i], queries[0])
	}
}

func TestPsqlParserHandleComments(t *testing.T) {
	r := require.New(t)

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

		r.Equal(sqlContent[i], queries[0])
	}
}

func TestPsqlParserHandleCommentsAndStringsMix(t *testing.T) {
	r := require.New(t)

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

		r.Equal(sqlContent[i], queries[0])
	}
}

func TestPsqlParserHandleTransactionBlock(t *testing.T) {
	r := require.New(t)

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

		r.Equal(sqlContent[i], queries[0])
	}
}

func TestPsqlParserFromSqlFile(t *testing.T) {
	r := require.New(t)

	sqlFilename := "queries_*.sql"
	sqlContent := "SELECT 1;"
	tempFile, _ := os.CreateTemp("", sqlFilename)
	tempFile.Write([]byte(sqlContent))

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	parser, _ := NewBuilder("psql").
		FromFile(tempFile.Name()).
		Build()
	queries := parser.Parse()

	r.Equal(1, len(queries))
}

func TestPsqlParserGCommand(t *testing.T) {
	r := require.New(t)

	sqlContent := `SELECT 1\g
SELECT 2\g result.txt
SELECT 3\g (format=unaligned tuples_only)
`
	parser, _ := NewBuilder("psql").
		WithContent(sqlContent).
		Build()

	queries := parser.Parse()

	r.Equal(3, len(queries))
	r.Equal("SELECT 1\\g\n", queries[0])
	r.Equal("SELECT 2\\g result.txt\n", queries[1])
	r.Equal("SELECT 3\\g (format=unaligned tuples_only)\n", queries[2])
}

func TestPsqlParserCrosstabviewCommand(t *testing.T) {
	r := require.New(t)

	sqlContent := `SELECT 1, 1, 1 \crosstabview
SELECT 2, 2, 2 \crosstabview
SELECT 3, 3, 3 \crosstabview
`
	parser, _ := NewBuilder("psql").
		WithContent(sqlContent).
		Build()
	queries := parser.Parse()

	r.Equal(3, len(queries))
}

func TestPsqlParserUnsupportedCommand(t *testing.T) {
	r := require.New(t)

	sqlContent := "SELECT 1\\unsupported\nSELECT 1\\\nSELECT 1;"
	parser, _ := NewBuilder("psql").
		WithContent(sqlContent).
		Build()
	queries := parser.Parse()

	r.Equal(1, len(queries))
}
