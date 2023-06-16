package dispatcher

import (
	"context"
	"fmt"

	"github.com/fljdin/dispatch/internal/models"
)

type DispatcherBuilder struct {
	dispatcher     *Dispatcher
	logfileName    string
	memorySize     int
	consoleEnabled bool
	err            error
}

func NewDispatcherBuilder() *DispatcherBuilder {
	ctx, cancel := context.WithCancel(context.Background())

	return &DispatcherBuilder{
		dispatcher: &Dispatcher{
			context: ctx,
			cancel:  cancel,
			workers: 1,
		},
		memorySize:     1,
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

func (db *DispatcherBuilder) WithMemorySize(size int) *DispatcherBuilder {
	db.memorySize = size
	return db
}

func (db *DispatcherBuilder) Build() (*Dispatcher, error) {
	if db.memorySize < 1 {
		db.err = fmt.Errorf("dispatcher need a positive memory size")
	}

	if db.dispatcher.workers < 1 {
		db.err = fmt.Errorf("dispatcher need a positive worker number")
	}

	db.dispatcher.tasks = make(chan models.Task, db.memorySize)
	db.dispatcher.results = make(chan models.TaskResult, db.memorySize)

	return db.dispatcher, db.err
}
