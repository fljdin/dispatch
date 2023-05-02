package dispatcher

import (
	"sync"

	"github.com/fljdin/dispatch/src/models"
)

type DispatcherMap struct {
	completed sync.Map
}

func (dm *DispatcherMap) Load(key int) int {
	status, ok := dm.completed.Load(key)

	if !ok {
		return models.Waiting
	}

	return status.(int)
}

func (dm *DispatcherMap) Store(key int, newStatus int) {
	var status int
	currentStatus := dm.Load(key)

	if currentStatus > newStatus {
		status = currentStatus
	} else {
		status = newStatus
	}

	dm.completed.Store(key, status)
}
