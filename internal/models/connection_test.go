package models_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestConnectionURIbyNameNotFound(t *testing.T) {
	var c Connections
	c = append(c, Connection{
		Name: "db",
		URI:  "postgresql://localhost",
	})
	_, err := c.GetURIByName("nowhere")
	require.Error(t, err)
}

var testConnections = []struct {
	name       string
	connection Connection
	expected   string
}{
	{
		name:       "TestConnectionCombinedURI1",
		connection: Connection{Host: "localhost", Port: 5432, Dbname: "postgres", User: "postgres", Password: "xxxxxxxx"},
		expected:   "postgresql://?dbname=postgres&host=localhost&password=xxxxxxxx&port=5432&user=postgres",
	},
	{
		name:       "TestConnectionCombinedURI2",
		connection: Connection{Host: "localhost", Port: 5432, Dbname: "postgres", User: "postgres"},
		expected:   "postgresql://?dbname=postgres&host=localhost&port=5432&user=postgres",
	},
	{
		name:       "TestConnectionCombinedURI3",
		connection: Connection{Dbname: "postgres", User: "postgres"},
		expected:   "postgresql://?dbname=postgres&user=postgres",
	},
}

func TestConnectionCombinedURI(t *testing.T) {
	for _, tc := range testConnections {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.connection.CombinedURI()
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestConnectionCombinedURIbyName(t *testing.T) {
	c := Connections{Connection{
		Name: "db",
		Host: "localhost",
	}}

	expected := "postgresql://?host=localhost"
	uri, _ := c.GetURIByName("db")

	assert.Equal(t, expected, uri)
}
