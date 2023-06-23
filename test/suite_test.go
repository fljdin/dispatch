package cmd_test

import (
	"context"
	"os"

	"github.com/fljdin/dispatch/internal/cmd"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

const (
	TestDatabase = "dispatch"
	TestRole     = "dispatch"
	TestPassword = "dispatch"
)

type Suite struct {
	suite.Suite
	Logfile  string
	Args     []string
	Output   func() string
	Expected func() string
}

func (suite *Suite) SetupSuite() {
	ctx := context.Background()
	conn, err := pgx.Connect(
		ctx, "postgresql://postgres:postgres@localhost:5432/postgres",
	)
	suite.Nil(err)

	setup, err := os.ReadFile("sql/setup.sql")
	suite.Nil(err)

	_, err = conn.Exec(ctx, string(setup))
	suite.Nil(err)
}

func (suite *Suite) SetupTest() {
	suite.Output = func() string {
		r := suite.Require()

		data, err := os.ReadFile(suite.Logfile)
		defer os.Remove(suite.Logfile)

		r.Nil(err)
		return string(data)
	}
	suite.Expected = func() string {
		r := suite.Require()

		data, err := os.ReadFile("expected/" + suite.Logfile)

		r.Nil(err)
		return string(data)
	}
}

func (suite *Suite) TestRun() {
	r := suite.Require()

	cmd.RootCmd.SetArgs(suite.Args)
	err := cmd.RootCmd.Execute()

	r.Nil(err)
	r.Equal(suite.Expected(), suite.Output())
}
