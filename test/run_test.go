package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

var s = &Suite{}

func TestTaskWithDefaultConnection(t *testing.T) {
	s.Args = []string{"run", "--config", "config/TaskWithDefaultConnection.yaml"}
	s.Logfile = "TaskWithDefaultConnection.log"

	suite.Run(t, s)
}

func TestWorkerForwardURIToGeneratedTasks(t *testing.T) {

}
