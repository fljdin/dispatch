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
	assert.Contains(t, err.Error(), "invalid type for parsing file")
}

func TestParserWithNullParseTYpe(t *testing.T) {
	pb, err := NewParserBuilder("").
		Build()

	require.NotNil(t, err)
	assert.Equal(t, "sh", pb.Type)
}
