package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionURIbyName(t *testing.T) {
	var c Connections
	given := "postgresql://localhost"
	c = append(c, Connection{
		Name: "db",
		URI:  given,
	})
	uri, _ := c.GetURIByName("db")

	assert.Equal(t, uri, given)
}
