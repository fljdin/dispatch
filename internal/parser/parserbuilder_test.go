package parser_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParserWithInvalidParseType(t *testing.T) {
	_, err := NewParserBuilder("unknown").
		Build()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "only psql type is supported")
}
