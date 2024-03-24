package dispatcher

import (
	"context"
	"fmt"

	"github.com/fljdin/dispatch/internal/queue"
	"github.com/fljdin/dispatch/internal/tasks"
)

type DispatcherBuilder struct {
	dispatcher Dispatcher
	err        error
}

func NewBuilder() *DispatcherBuilder {
	ctx, cancel := context.WithCancel(context.Background())

	return &DispatcherBuilder{
		dispatcher: Dispatcher{
			context:   ctx,
			cancel:    cancel,
			processes: 1,
		},
	}
}

func (db *DispatcherBuilder) WithProcesses(count int) *DispatcherBuilder {
	db.dispatcher.processes = count
	return db
}

func (db *DispatcherBuilder) Build() (Dispatcher, error) {
	if db.dispatcher.processes < 1 {
		db.err = fmt.Errorf("dispatcher need a positive processes number")
	}

	db.dispatcher.memory = &Memory{
		queue:   queue.New(),
		results: make(chan tasks.Result, 10),
	}

	db.dispatcher.monitor = NewMonitor(
		db.dispatcher.memory,
		db.dispatcher.context,
	)

	return db.dispatcher, db.err
}
