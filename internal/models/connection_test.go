package models_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestConnectionURIbyName(t *testing.T) {
	var c Connections
	expected := "postgresql://localhost"
	c = append(c, Connection{
		Name: "db",
		URI:  expected,
	})
	uri, _ := c.GetURIByName("db")

	assert.Equal(t, expected, uri)
}
