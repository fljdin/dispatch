package dispatcher_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/dispatcher"
	"github.com/stretchr/testify/require"
)

func TestDispatcherBuilderWithNegativeProcsNumber(t *testing.T) {
	r := require.New(t)

	_, err := NewBuilder().
		WithProcesses(0).
		Build()

	r.NotNil(err)
	r.Contains(err.Error(), "need a positive processes number")
}
