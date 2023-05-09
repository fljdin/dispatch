package dispatcher_test

import (
	"context"
	"testing"

	. "github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/stretchr/testify/assert"
)

func TestDispatcherBuilderWithInvalidMemorySize(t *testing.T) {
	_, err := NewDispatcherBuilder(context.Background()).
		WithMemorySize(0).
		Build()

	if assert.NotEqual(t, nil, err) {
		assert.Contains(t, err.Error(), "need a positive memory size")
	}
}

func TestDispatcherBuilderWithInvalidTraceFile(t *testing.T) {
	_, err := NewDispatcherBuilder(context.Background()).
		WithTraceFile("not/exists.out").
		Build()

	if assert.NotEqual(t, nil, err) {
		assert.Contains(t, err.Error(), "no such file or directory")
	}
}

func TestDispatcherBuilderWithNegativeWorkerNumber(t *testing.T) {
	_, err := NewDispatcherBuilder(context.Background()).
		WithWorkerNumber(0).
		Build()

	if assert.NotEqual(t, nil, err) {
		assert.Contains(t, err.Error(), "need a positive worker number")
	}
}
