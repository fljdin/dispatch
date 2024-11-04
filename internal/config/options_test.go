package config_test

import (
	"runtime"
	"testing"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/knadh/koanf/v2"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
)

func TestLoadFlagFromShortOption(t *testing.T) {
	r := require.New(t)
	k := koanf.New(".")

	opts := config.Flags()
	opts.Parse([]string{
		"-c", "config.yaml",
		"-P", "2",
		"-o", "/dev/stdout",
		"-v",
	})
	r.NoError(config.LoadFlags(k, opts))

	r.Equal("config.yaml", k.String("config"))
	r.Equal(2, k.Int("procs"))
	r.Equal("/dev/stdout", k.String("output"))
	r.True(k.Bool("verbose"))
}

func TestLoadFlagPrecedenceOverYAML(t *testing.T) {
	r := require.New(t)
	k := koanf.New(".")

	opts := config.Flags()
	opts.Parse([]string{
		"--config", "config.yaml",
		"--procs", "4",
		"--output", "/dev/stdout",
		"--verbose",
	})
	r.NoError(config.LoadFlags(k, opts))

	yaml := []byte(dedent.Dedent(`
		procs: 1
		verbose: false
		output: /dev/null
	`))
	r.NoError(config.LoadYamlRaw(k, yaml))

	r.Equal(4, k.Int("procs"))
	r.Equal("/dev/stdout", k.String("output"))
	r.True(k.Bool("verbose"))
}

func TestLoadConfigIsRequired(t *testing.T) {
	r := require.New(t)
	k := koanf.New(".")

	opts := config.Flags()
	opts.Parse([]string{})

	r.Error(config.LoadFlags(k, opts))
}

func TestLoadVersionOnly(t *testing.T) {
	r := require.New(t)
	k := koanf.New(".")

	opts := config.Flags()
	opts.Parse([]string{"-version"})
	r.NoError(config.LoadFlags(k, opts))
}

func TestLoadProcessFromYAML(t *testing.T) {
	r := require.New(t)
	k := koanf.New(".")

	opts := config.Flags()
	opts.Parse([]string{
		"-c", "config.yaml",
	})
	r.NoError(config.LoadFlags(k, opts))

	yaml := []byte(dedent.Dedent(`
		procs: 4
	`))
	r.NoError(config.LoadYamlRaw(k, yaml))

	r.Equal(4, k.Int("procs"))
}

func TestLoadDefaultProcessNumber(t *testing.T) {
	r := require.New(t)
	k := koanf.New(".")

	opts := config.Flags()
	opts.Parse([]string{
		"-c", "config.yaml",
	})
	r.NoError(config.LoadFlags(k, opts))

	procs := config.ValidateProcs(k.Int("procs"))
	r.Equal(config.ProcessesDefault, procs)
}

func TestValidateProcessNumberBoundary(t *testing.T) {
	r := require.New(t)

	procs := config.ValidateProcs(-1)
	r.Equal(config.ProcessesDefault, procs)

	procs = config.ValidateProcs(64)
	r.Equal(runtime.NumCPU(), procs)
}
