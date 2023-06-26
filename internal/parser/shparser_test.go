package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShParserWithContent(t *testing.T) {
	r := require.New(t)

	shContent := "true\nfalse"
	parser, _ := NewBuilder("sh").
		WithContent(shContent).
		Build()
	commands := parser.Parse()

	r.Equal(2, len(commands))
	r.Equal("true", commands[0])
	r.Equal("false", commands[1])
}

func TestShParserIgnoreEmptyLines(t *testing.T) {
	r := require.New(t)

	shContent := "\n\n\ntrue\n\n"
	parser, _ := NewBuilder("sh").
		WithContent(shContent).
		Build()
	commands := parser.Parse()

	r.Equal(1, len(commands))
	r.Equal("true", commands[0])
}

func TestShParserWithEscape(t *testing.T) {
	r := require.New(t)

	shContent := "true\\\n && false"
	parser, _ := NewBuilder("sh").
		WithContent(shContent).
		Build()
	commands := parser.Parse()

	r.Equal(1, len(commands))
	r.Equal("true && false", commands[0])
}

func TestShParserWithComment(t *testing.T) {
	r := require.New(t)

	shContent := "true # \\\n false"
	parser, _ := NewBuilder("sh").
		WithContent(shContent).
		Build()
	commands := parser.Parse()

	r.Equal(2, len(commands))
	r.Equal("true", commands[0])
	r.Equal("false", commands[1])
}
