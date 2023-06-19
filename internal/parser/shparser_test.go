package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShParserWithContent(t *testing.T) {
	shContent := `true\nfalse`

	parser, _ := NewParserBuilder("sh").
		WithContent(shContent).
		Build()

	commands := parser.Parse()
	assert.Equal(t, 2, len(commands))
}
