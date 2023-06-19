package parser_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserWithInvalidParseType(t *testing.T) {
	_, err := NewBuilder("unknown").
		Build()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "is not supported")
}
