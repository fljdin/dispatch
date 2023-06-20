package dispatcher_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDispatcherBuilderWithInvalidTraceFile(t *testing.T) {
	_, err := NewDispatcherBuilder().
		WithLogfile("not/exists.out").
		Build()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestDispatcherBuilderWithNegativeWorkerNumber(t *testing.T) {
	_, err := NewDispatcherBuilder().
		WithWorkerNumber(0).
		Build()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "need a positive worker number")
}
