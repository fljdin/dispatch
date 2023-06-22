package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHelloWorld(t *testing.T) {
	suite.Run(t, &Suite{
		Args:    []string{"run", "--config", "config/HelloWorld.yaml"},
		Logfile: "HelloWorld.log",
	})
}

func TestWorkerForwardURIToGeneratedTasks(t *testing.T) {

}
