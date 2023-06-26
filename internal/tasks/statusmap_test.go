package tasks_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/require"
)

func TestStatusMapLoad(t *testing.T) {
	r := require.New(t)

	sm := StatusMap{}
	status := sm.Get(1)

	r.Equal(Waiting, status)
}

func TestStatusMapStore(t *testing.T) {
	r := require.New(t)

	sm := StatusMap{}
	expected := Succeeded

	sm.Set(1, Succeeded)
	status := sm.Get(1)

	r.Equal(expected, status)
}

func TestStatusMapUpdate(t *testing.T) {
	r := require.New(t)

	sm := StatusMap{}

	expected := Succeeded
	sm.Set(1, Waiting)
	sm.Set(1, Succeeded)
	status := sm.Get(1)

	r.Equal(expected, status)

	expected = Failed
	sm.Set(2, Failed)
	sm.Set(2, Succeeded)
	status = sm.Get(2)

	r.Equal(expected, status)
}
