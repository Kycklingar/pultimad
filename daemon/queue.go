package daemon

import "sync"

type queue struct {
	tasks []task
	l     *sync.Mutex
}

func (q *queue) push(t task) {
	q.l.Lock()
	defer q.l.Unlock()
	q.tasks = append(q.tasks, t)
}

func (q *queue) pop() task {
	q.l.Lock()
	defer q.l.Unlock()
	if len(q.tasks) == 0 {
		return task{}
	}

	t := q.tasks[0]
	q.tasks = q.tasks[1:]
	return t
}
