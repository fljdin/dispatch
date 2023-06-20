package tasks

import (
	"sync"
)

const (
	Waiting int = iota
	Ready
	Succeeded
	Failed
	Interrupted
)

type StatusMap struct {
	sm sync.Map
}

func (dm *StatusMap) Get(id int) int {
	status, exists := dm.sm.Load(id)

	if !exists {
		return Waiting
	}

	return status.(int)
}

func (dm *StatusMap) Set(id int, newStatus int) {
	var status int
	currentStatus := dm.Get(id)

	if currentStatus > newStatus {
		status = currentStatus
	} else {
		status = newStatus
	}

	dm.sm.Store(id, status)
}
