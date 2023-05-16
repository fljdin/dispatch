package models

import (
	"sync"
)

type StatusMap struct {
	statuses sync.Map
}

func (dm *StatusMap) Load(taskID int) int {
	status, exists := dm.statuses.Load(taskID)

	if !exists {
		return Waiting
	}

	return status.(int)
}

func (dm *StatusMap) Store(taskID int, newStatus int) {
	var status int
	currentStatus := dm.Load(taskID)

	if currentStatus > newStatus {
		status = currentStatus
	} else {
		status = newStatus
	}

	dm.statuses.Store(taskID, status)
}
