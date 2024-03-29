package config

import (
	"flag"
	"fmt"
	"runtime"
	"strconv"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/basicflag"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

func Flags() *flag.FlagSet {
	f := flag.NewFlagSet("", flag.ExitOnError)

	f.Bool("v", false, "")
	f.Bool("verbose", false, "")
	f.Bool("version", false, "")
	f.Int("P", 0, "")
	f.Int("procs", 0, "")
	f.String("c", "", "")
	f.String("config", "", "")
	f.String("o", "", "")
	f.String("output", "", "")

	return f
}

func LoadFlags(k *koanf.Koanf, opts *flag.FlagSet) error {
	return k.Load(
		basicflag.Provider(opts, "koanf"), nil,
		MergeFunc,
	)
}

func LoadYaml(k *koanf.Koanf, path string) error {
	return k.Load(
		file.Provider(path), yaml.Parser(),
		MergeFunc,
	)
}

func LoadYamlRaw(k *koanf.Koanf, raw []byte) error {
	return k.Load(
		rawbytes.Provider(raw), yaml.Parser(),
		MergeFunc,
	)
}

// boundary check for the number of processes
func ValidateProcs(procs int) int {
	if procs < 1 {
		return ProcessesDefault
	}

	if procs > runtime.NumCPU() {
		return runtime.NumCPU()
	}

	return procs
}

var MergeFunc koanf.Option = koanf.WithMergeFunc(func(src, dest map[string]any) error {
	var IsZero = func(v any) bool {
		return v == 0 || v == "0" ||
			v == "" || v == "false"
	}

	var IsDefined = func(v any) bool {
		return v != nil && !IsZero(v)
	}

	for k, v := range src {
		// do not overwrite with zero values
		if IsZero(v) {
			continue
		}

		// do not overwrite an already defined value
		if IsDefined(dest[k]) {
			continue
		}

		switch k {
		case "c", "config":
			dest["config"] = v

		case "o", "output":
			dest["output"] = v

		case "v", "verbose":
			dest["verbose"] = v

		case "P", "procs":
			switch v := v.(type) {
			case int:
				// when value comes from yaml, it's an int
				dest["procs"] = v

			case string:
				// when value comes from flag, it's a string
				dest["procs"], _ = strconv.Atoi(v)
			}

		default:
			dest[k] = v
		}
	}

	// stop here if the version flag is set
	if IsDefined(dest["version"]) {
		return nil
	}

	// raise an error if the config file is missing
	if !IsDefined(dest["config"]) {
		return fmt.Errorf("missing configuration file")
	}

	return nil
})
