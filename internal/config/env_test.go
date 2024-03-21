package config_test

import (
	"testing"

	. "github.com/fljdin/dispatch/internal/config"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

var (
	yml string = dedent.Dedent(`
		name: local
		variables:
		  PGHOST: localhost
		  PGDATABASE: db
	`)

	envs Environments = Environments{
		Environment{
			Name: "local",
			Variables: map[string]string{
				"PGHOST": "localhost",
				"PGUSER": "db",
			},
		},
	}
)

func TestEnvironmentUnmarshal(t *testing.T) {
	r := require.New(t)

	var env Environment
	yaml.Unmarshal([]byte(yml), &env)

	r.Equal("localhost", env.Variables["PGHOST"])
}

func TestEnvironmentByName(t *testing.T) {
	r := require.New(t)

	env, _ := envs.ByName("local")
	r.Equal("localhost", env.Variables["PGHOST"])
}

func TestEnvironmentNotFound(t *testing.T) {
	r := require.New(t)

	_, err := envs.ByName("remote")
	r.Error(err)
}
