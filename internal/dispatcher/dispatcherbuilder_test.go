package dispatcher_test

import (
	"context"
	"testing"

	. "github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDispatcherBuilderWithInvalidMemorySize(t *testing.T) {
	_, err := NewDispatcherBuilder(context.Background()).
		WithMemorySize(0).
		Build()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "need a positive memory size")
}

func TestDispatcherBuilderWithInvalidTraceFile(t *testing.T) {
	_, err := NewDispatcherBuilder(context.Background()).
		WithLogfile("not/exists.out").
		Build()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestDispatcherBuilderWithNegativeWorkerNumber(t *testing.T) {
	_, err := NewDispatcherBuilder(context.Background()).
		WithWorkerNumber(0).
		Build()

	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "need a positive worker number")
}
