package queue

import (
	"sync"

	"github.com/fljdin/dispatch/internal/status"
)

type StatusMap struct {
	m   map[int][]status.Status
	mut sync.Mutex
}

func NewStatusMap() *StatusMap {
	return &StatusMap{
		m: make(map[int][]status.Status),
	}
}

func (sm *StatusMap) Get(id int) status.Status {
	sm.mut.Lock()
	defer sm.mut.Unlock()

	values, exists := sm.m[id]
	if !exists {
		return status.Waiting
	}

	for _, v := range values {
		if v == status.Interrupted {
			return status.Interrupted
		}
		if v == status.Failed {
			return status.Failed
		}
		if v == status.Waiting {
			return status.Waiting
		}
	}
	return status.Succeeded
}

func (sm *StatusMap) Set(id, subid int, status status.Status) {
	sm.mut.Lock()
	defer sm.mut.Unlock()

	values := sm.m[id]

	if subid >= len(values) {
		values = append(values, status)
	} else {
		values[subid] = status
	}

	sm.m[id] = values
}
