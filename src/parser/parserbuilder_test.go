package parser_test

import (
	"testing"

	. "github.com/fljdin/dispatch/src/parser"
	"github.com/stretchr/testify/assert"
)

func TestParserFromInvalidParseType(t *testing.T) {
	_, err := NewParserBuilder("unknown").
		Build()

	if assert.NotEqual(t, nil, err) {
		assert.Contains(t, err.Error(), "invalid type for parsing file")
	}
}
