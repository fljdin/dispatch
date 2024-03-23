package dispatcher

import (
	"context"
	"fmt"

	"github.com/fljdin/dispatch/internal/queue"
	"github.com/fljdin/dispatch/internal/tasks"
)

type DispatcherBuilder struct {
	dispatcher     Dispatcher
	logfileName    string
	consoleEnabled bool
	err            error
}

func NewBuilder() *DispatcherBuilder {
	ctx, cancel := context.WithCancel(context.Background())

	return &DispatcherBuilder{
		dispatcher: Dispatcher{
			context:   ctx,
			cancel:    cancel,
			processes: 1,
		},
		consoleEnabled: false,
	}
}

func (db *DispatcherBuilder) WithLogfile(filename string) *DispatcherBuilder {
	db.logfileName = filename
	return db
}

func (db *DispatcherBuilder) WithConsole() *DispatcherBuilder {
	db.consoleEnabled = true
	return db
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

	if db.consoleEnabled {
		db.dispatcher.monitor.WithConsole()
	}

	if err := db.dispatcher.monitor.WithTrace(db.logfileName); err != nil {
		db.err = err
	}

	return db.dispatcher, db.err
}
