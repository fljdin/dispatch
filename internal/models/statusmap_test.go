package models_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestStatusMapLoad(t *testing.T) {
	sm := StatusMap{}

	status := sm.Load(1)

	assert.Equal(t, Waiting, status)
}

func TestStatusMapStore(t *testing.T) {
	sm := StatusMap{}
	expected := Succeeded

	sm.Store(1, Succeeded)
	status := sm.Load(1)

	assert.Equal(t, expected, status)
}

func TestStatusMapUpdate(t *testing.T) {
	sm := StatusMap{}

	expected := Succeeded
	sm.Store(1, Waiting)
	sm.Store(1, Succeeded)
	status := sm.Load(1)

	assert.Equal(t, expected, status)

	expected = Failed
	sm.Store(2, Failed)
	sm.Store(2, Succeeded)
	status = sm.Load(2)

	assert.Equal(t, expected, status)

}
