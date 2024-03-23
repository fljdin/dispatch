package dispatcher_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/stretchr/testify/require"
)

func TestDispatcherBuilderWithInvalidTraceFile(t *testing.T) {
	r := require.New(t)

	_, err := NewBuilder().
		WithLogfile("not/exists.out").
		Build()

	r.NotNil(err)
	r.Contains(err.Error(), "no such file or directory")
}

func TestDispatcherBuilderWithNegativeProcsNumber(t *testing.T) {
	r := require.New(t)

	_, err := NewBuilder().
		WithProcesses(0).
		Build()

	r.NotNil(err)
	r.Contains(err.Error(), "need a positive processes number")
}
