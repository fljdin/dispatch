package tasks_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/assert"
)

func TestStatusMapLoad(t *testing.T) {
	sm := StatusMap{}

	status := sm.Get(1)

	assert.Equal(t, Waiting, status)
}

func TestStatusMapStore(t *testing.T) {
	sm := StatusMap{}
	expected := Succeeded

	sm.Set(1, Succeeded)
	status := sm.Get(1)

	assert.Equal(t, expected, status)
}

func TestStatusMapUpdate(t *testing.T) {
	sm := StatusMap{}

	expected := Succeeded
	sm.Set(1, Waiting)
	sm.Set(1, Succeeded)
	status := sm.Get(1)

	assert.Equal(t, expected, status)

	expected = Failed
	sm.Set(2, Failed)
	sm.Set(2, Succeeded)
	status = sm.Get(2)

	assert.Equal(t, expected, status)
}
