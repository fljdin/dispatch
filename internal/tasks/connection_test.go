package tasks_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/tasks"
	"github.com/stretchr/testify/require"
)

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
	{
		name:       "TestConnectionCombinedURI4",
		connection: Connection{Service: "mydb", User: "postgres"},
		expected:   "postgresql://?service=mydb&user=postgres",
	},
}

func TestConnectionCombinedURI(t *testing.T) {
	r := require.New(t)

	for _, tc := range testConnections {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.connection.CombinedURI()
			r.Equal(tc.expected, actual)
		})
	}
}

func TestConnectionsURIbyName(t *testing.T) {
	r := require.New(t)

	var c Connections
	expected := "postgresql://localhost"
	c = append(c, Connection{
		Name: "db",
		URI:  expected,
	})
	uri, _ := c.GetURIByName("db")

	r.Equal(expected, uri)
}

func TestConnectionsURIbyNameNotFound(t *testing.T) {
	r := require.New(t)

	var c Connections
	c = append(c, Connection{
		Name: "db",
		URI:  "postgresql://localhost",
	})
	_, err := c.GetURIByName("nowhere")

	r.NotNil(err)
}

func TestConnectionsCombinedURIbyName(t *testing.T) {
	r := require.New(t)

	c := Connections{Connection{
		Name: "db",
		Host: "localhost",
	}}
	expected := "postgresql://?host=localhost"
	uri, _ := c.GetURIByName("db")

	r.Equal(expected, uri)
}
