package dispatcher

import (
	"context"
	"fmt"

	"github.com/fljdin/dispatch/internal/models"
)

type DispatcherBuilder struct {
	dispatcher     Dispatcher
	traceFilename  string
	memorySize     int
	consoleEnabled bool
	err            error
}

func NewDispatcherBuilder(ctx context.Context) *DispatcherBuilder {
	ctx, cancel := context.WithCancel(ctx)

	return &DispatcherBuilder{
		dispatcher: Dispatcher{
			context: ctx,
			cancel:  cancel,
			workers: 1,
		},
		memorySize:     1,
		consoleEnabled: false,
	}
}

func (db *DispatcherBuilder) WithTraceFile(filename string) *DispatcherBuilder {
	db.traceFilename = filename
	return db
}

func (db *DispatcherBuilder) WithConsole() *DispatcherBuilder {
	db.consoleEnabled = true
	return db
}

func (db *DispatcherBuilder) WithWorkerNumber(count int) *DispatcherBuilder {
	db.dispatcher.workers = count
	return db
}

func (db *DispatcherBuilder) WithMemorySize(size int) *DispatcherBuilder {
	db.memorySize = size
	return db
}

func (db *DispatcherBuilder) Build() (Dispatcher, error) {
	if db.memorySize < 1 {
		db.err = fmt.Errorf("dispatcher need a positive memory size")
	}

	if db.dispatcher.workers < 1 {
		db.err = fmt.Errorf("dispatcher need a positive worker number")
	}

	db.dispatcher.memory = &SharedMemory{
		tasks:   make(chan models.Task, db.memorySize),
		results: make(chan models.TaskResult, db.memorySize),
	}

	db.dispatcher.observer = &Observer{
		memory:  db.dispatcher.memory,
		context: db.dispatcher.context,
	}

	if err := db.dispatcher.observer.WithTrace(db.traceFilename); err != nil {
		db.err = err
	}

	if db.consoleEnabled {
		db.dispatcher.observer.WithConsole()
	}

	return db.dispatcher, db.err
}
