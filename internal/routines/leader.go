package routines

import (
	"context"

	"github.com/fljdin/dispatch/internal/config"
)

type Leader struct {
	cancel  func()
	context context.Context
	*Memory
}

func NewLeader(procs int) Leader {
	ctx, cancel := context.WithCancel(context.Background())

	leader := Leader{
		context: ctx,
		cancel:  cancel,
		Memory: &Memory{
			processes: procs,
			queue:     NewQueue(),
			tasks:     make(chan config.Task, procs),
		},
	}

	return leader
}

func (l *Leader) Wait() {
	l.FillTasks()
	l.launchProcesses()
	l.WaitForTasks()
	l.cancel()
	l.WaitForProcesses()
}

func (l Leader) launchProcesses() {
	for i := 1; i <= l.processes; i++ {
		process := Process{
			memory:  l.Memory,
			context: l.context,
		}
		go process.Start()
	}
}
