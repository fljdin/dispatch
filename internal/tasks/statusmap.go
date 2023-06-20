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
	m sync.Map
}

func (sm *StatusMap) Get(id int) int {
	status, exists := sm.m.Load(id)

	if !exists {
		return Waiting
	}

	return status.(int)
}

func (sm *StatusMap) Set(id int, newStatus int) {
	var status int
	currentStatus := sm.Get(id)

	if currentStatus > newStatus {
		status = currentStatus
	} else {
		status = newStatus
	}

	sm.m.Store(id, status)
}
