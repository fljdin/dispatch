package parser_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/parser"
	"github.com/stretchr/testify/require"
)

func TestParserWithInvalidParseType(t *testing.T) {
	r := require.New(t)

	_, err := NewBuilder("unknown").
		Build()

	r.NotNil(err)
	r.Contains(err.Error(), "is not supported")
}
