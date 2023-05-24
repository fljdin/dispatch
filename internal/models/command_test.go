package models_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCommandBasicRun(t *testing.T) {
	cmd := Command{
		Type:    "sh",
		Command: "echo test",
	}
	result := cmd.Run()

	assert.Equal(t, Succeeded, result.Status)
	assert.Contains(t, result.Output, "test")
}
