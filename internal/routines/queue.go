package routines

import (
	"sync"

	"github.com/fljdin/dispatch/internal/config"
	"github.com/fljdin/dispatch/internal/status"
	om "github.com/wk8/go-ordered-map/v2"
)

type Queue struct {
	mut   sync.Mutex
	tasks *om.OrderedMap[int, []config.Task]
}

func NewQueue() Queue {
	return Queue{
		tasks: om.New[int, []config.Task](),
	}
}

func (q *Queue) Add(t config.Task) {
	q.mut.Lock()
	defer q.mut.Unlock()

	subs, exists := q.tasks.Get(t.Identifier.ID)
	if !exists {
		q.tasks.Set(t.Identifier.ID, []config.Task{t})
		return
	}

	subs = append(subs, t)
	q.tasks.Set(t.Identifier.ID, subs)
}

func (q *Queue) Update(tid config.TaskIdentifier, s status.Status) {
	q.mut.Lock()
	defer q.mut.Unlock()

	// update status in queue
	if t, exists := q.tasks.Get(tid.ID); exists {
		t[tid.SubID].Status = s
		q.tasks.Set(tid.ID, t)
		return
	}
}

func (q *Queue) Next() (config.Task, bool) {
	q.mut.Lock()
	defer q.mut.Unlock()

	for pair := q.tasks.Oldest(); pair != nil; pair = pair.Next() {
	next:
		for _, sub := range pair.Value {
			// ignore already processed tasks
			if sub.Status != status.Waiting {
				continue next
			}

			for _, id := range sub.Depends {
				dep := q.Evaluate(id)

				// interrupt if a dependency failed
				if dep.IsFailed() {
					sub.Status = status.Interrupted
					return sub, true
				}

				// wait until a dependency succeeds
				if !dep.IsSucceeded() {
					continue next
				}
			}

			sub.Status = status.Ready
			return sub, true
		}
	}

	// no task found
	return config.Task{}, false
}

func (q *Queue) Evaluate(id int) status.Status {
	subs, exists := q.tasks.Get(id)
	if !exists {
		return status.Waiting
	}

	for _, sub := range subs {
		switch sub.Status {
		case status.Waiting, status.Ready, status.Running:
			return status.Waiting

		case status.Interrupted:
			return status.Interrupted

		case status.Failed:
			return status.Failed

		default:
			continue
		}
	}

	return status.Succeeded
}
