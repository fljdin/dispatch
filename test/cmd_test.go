package cmd_test

import (
	"os"

	"github.com/fljdin/dispatch/internal/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	Logfile  string
	Args     []string
	Output   func() string
	Expected func() string
}

func (suite *Suite) SetupTest() {
	suite.Output = func() string {
		data, err := os.ReadFile(suite.Logfile)
		defer os.Remove(suite.Logfile)
		require.Nil(suite.T(), err)
		return string(data)
	}
	suite.Expected = func() string {
		data, err := os.ReadFile("expected/" + suite.Logfile)
		require.Nil(suite.T(), err)
		return string(data)
	}
}

func (suite *Suite) Run() {
	cmd.RootCmd.SetArgs(suite.Args)
	err := cmd.RootCmd.Execute()

	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.Expected(), suite.Output())
}
