package queue

import (
	"sync"

	. "github.com/fljdin/dispatch/internal/tasks"
)

type Status map[int]int

func (s Status) Get() (status int) {
	for _, v := range s {
		if v == Interrupted {
			return Interrupted
		}
		if v == Failed {
			return Failed
		}
		if v == Waiting {
			return Waiting
		}
	}
	return Succeeded
}

type StatusMap struct {
	m   map[int]Status
	mut sync.Mutex
}

func NewStatusMap() *StatusMap {
	return &StatusMap{
		m: make(map[int]Status),
	}
}

func (sm *StatusMap) Get(id int) int {
	sm.mut.Lock()
	defer sm.mut.Unlock()

	values, exists := sm.m[id]
	if !exists {
		return Waiting
	}

	return values.Get()
}

func (sm *StatusMap) Set(id, subid, status int) {
	sm.mut.Lock()
	defer sm.mut.Unlock()

	values, exists := sm.m[id]
	if !exists {
		values = make(Status)
	}

	values[subid] = status
	sm.m[id] = values
}
