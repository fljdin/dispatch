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

func TestDispatcherBuilderWithNegativeWorkerNumber(t *testing.T) {
	r := require.New(t)

	_, err := NewBuilder().
		WithWorkerNumber(0).
		Build()

	r.NotNil(err)
	r.Contains(err.Error(), "need a positive worker number")
}
