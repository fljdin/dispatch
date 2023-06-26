package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShParserWithContent(t *testing.T) {
	shContent := `true
false`

	parser, _ := NewBuilder("sh").
		WithContent(shContent).
		Build()
	commands := parser.Parse()

	assert.Equal(t, 2, len(commands))
	assert.Equal(t, "true", commands[0])
	assert.Equal(t, "false", commands[1])
}
