package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteWithoutConfiguration(t *testing.T) {
	os.Args = []string{"dispatch"}
	rc := runExecute()

	assert.NotEqual(t, rc, 0)
}

func TestExecuteWithInvalidYAML(t *testing.T) {
	yamlFile := "test.yaml"
	yamlContent := ""
	tempFile, _ := ioutil.TempFile(".", yamlFile)

	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	tempFile.Write([]byte(yamlContent))

	os.Args = []string{"dispatch", "-c", yamlFile}
	rc := runExecute()

	assert.NotEqual(t, rc, 0)
}

func runExecute() int {

	// restore previous args, stdout and stderr at end
	args := os.Args
	stdout := os.Stdout
	stderr := os.Stderr

	defer func() {
		os.Args = args
		os.Stdout = stdout
		os.Stderr = stderr
	}()

	// capture stdout and stderr
	_, wout, _ := os.Pipe()
	_, werr, _ := os.Pipe()

	os.Stdout = wout
	os.Stderr = werr

	return execute()
}
