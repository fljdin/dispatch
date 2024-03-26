package queue_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/queue"
	. "github.com/fljdin/dispatch/internal/status"
	"github.com/stretchr/testify/require"
)

func TestStatusMapLoad(t *testing.T) {
	r := require.New(t)

	sm := NewStatusMap()
	status := sm.Get(1)

	r.Equal(Waiting, status)
}

func TestStatusMapStore(t *testing.T) {
	r := require.New(t)

	sm := NewStatusMap()
	sm.Set(1, 0, Succeeded)

	status := sm.Get(1)
	r.Equal(Succeeded, status)
}

func TestStatusMapUpdate(t *testing.T) {
	r := require.New(t)
	sm := NewStatusMap()

	sm.Set(1, 0, Waiting)
	sm.Set(1, 0, Succeeded)
	status := sm.Get(1)

	r.Equal(Succeeded, status)

}

func TestStatusMapWaiting(t *testing.T) {
	r := require.New(t)
	sm := NewStatusMap()

	sm.Set(1, 0, Ready)
	sm.Set(1, 1, Waiting)
	sm.Set(1, 2, Succeeded)

	status := sm.Get(1)
	r.Equal(Waiting, status)

	sm.Set(1, 1, Succeeded)

	status = sm.Get(1)
	r.Equal(Succeeded, status)
}

func TestStatusMapFailed(t *testing.T) {
	r := require.New(t)
	sm := NewStatusMap()

	sm.Set(2, 0, Ready)
	sm.Set(2, 0, Failed)
	sm.Set(2, 1, Succeeded)

	status := sm.Get(2)
	r.Equal(Failed, status)
}

func TestStatusMapInterrupted(t *testing.T) {
	r := require.New(t)
	sm := NewStatusMap()

	sm.Set(3, 0, Ready)
	sm.Set(3, 0, Interrupted)
	sm.Set(3, 1, Succeeded)

	status := sm.Get(3)
	r.Equal(Interrupted, status)
}
