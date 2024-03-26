package queue

import (
	"sync"

	"github.com/fljdin/dispatch/internal/status"
	"github.com/fljdin/dispatch/internal/tasks"
	om "github.com/wk8/go-ordered-map/v2"
)

type Queue struct {
	mut   sync.Mutex
	tasks *om.OrderedMap[int, []tasks.Task]
}

func New() Queue {
	return Queue{
		tasks: om.New[int, []tasks.Task](),
	}
}

func (q *Queue) Add(t tasks.Task) {
	q.mut.Lock()
	defer q.mut.Unlock()

	subs, exists := q.tasks.Get(t.ID)
	if !exists {
		q.tasks.Set(t.ID, []tasks.Task{t})
		return
	}

	subs = append(subs, t)
	q.tasks.Set(t.ID, subs)
}

func (q *Queue) Update(id, subid int, s status.Status) {
	q.mut.Lock()
	defer q.mut.Unlock()

	// update status in queue
	if t, exists := q.tasks.Get(id); exists {
		t[subid].Status = s
		q.tasks.Set(id, t)
		return
	}
}

func (q *Queue) Next() (tasks.Task, bool) {
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
	return tasks.Task{}, false
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

		case status.Failed, status.Interrupted:
			return status.Failed

		default:
			continue
		}
	}

	return status.Succeeded
}
