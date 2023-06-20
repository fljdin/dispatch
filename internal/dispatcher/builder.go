package dispatcher

import (
	"context"
	"fmt"

	"github.com/fljdin/dispatch/internal/task"
)

type DispatcherBuilder struct {
	dispatcher     Dispatcher
	logfileName    string
	consoleEnabled bool
	err            error
}

func NewDispatcherBuilder() *DispatcherBuilder {
	ctx, cancel := context.WithCancel(context.Background())

	return &DispatcherBuilder{
		dispatcher: Dispatcher{
			context: ctx,
			cancel:  cancel,
			workers: 1,
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

func (db *DispatcherBuilder) WithWorkerNumber(count int) *DispatcherBuilder {
	db.dispatcher.workers = count
	return db
}

func (db *DispatcherBuilder) Build() (Dispatcher, error) {
	if db.dispatcher.workers < 1 {
		db.err = fmt.Errorf("dispatcher need a positive worker number")
	}

	db.dispatcher.memory = &Memory{
		queue:   task.NewQueue(),
		results: make(chan task.Result, 10),
	}

	db.dispatcher.observer = &Observer{
		memory:  db.dispatcher.memory,
		context: db.dispatcher.context,
	}

	if err := db.dispatcher.observer.WithTrace(db.logfileName); err != nil {
		db.err = err
	}

	if db.consoleEnabled {
		db.dispatcher.observer.WithConsole()
	}

	return db.dispatcher, db.err
}
