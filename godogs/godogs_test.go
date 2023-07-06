package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/fljdin/dispatch/internal/cmd"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/cucumber/godog"
)

// godogsCtxKey is the key used to store the available godogs in the context.Context.
type godogsCtxKey struct{}

var configurationFilePath string
var expectedLogFilePath string
var generatedLogFilePath string

func aConfigurationFilein(ctx context.Context, filepath string) (context.Context, error) {
	configurationFilePath := filepath
}

func iExecuteDispatcher(ctx context.Context) (context.Context, error) {
	cmd.RootCmd.SetArgs([]string{"run", "--config", configurationFilePath})
	err := cmd.RootCmd.Execute()
	require.Nil(t, err)
}

func thereShouldBeRemaining(ctx context.Context, remaining int) error {
	available, ok := ctx.Value(godogsCtxKey{}).(int)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if available != remaining {
		return fmt.Errorf("expected %d godogs to be remaining, but there is %d", remaining, available)
	}

	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	sc.Step(`^I eat (\d+)$`, iEat)
	sc.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
}
