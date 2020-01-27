package daemon

import "sync"

type queue struct {
	tasks []Taskif
	l     *sync.Mutex
}

func newQueue() *queue {
	return &queue{l: &sync.Mutex{}}
}

func (q *queue) push(t []Taskif) {
	q.l.Lock()
	defer q.l.Unlock()
	q.tasks = append(q.tasks, t...)
}

func (q *queue) pop() Taskif {
	q.l.Lock()
	defer q.l.Unlock()
	if len(q.tasks) == 0 {
		return nil
	}

	t := q.tasks[0]
	q.tasks = q.tasks[1:]
	return t
}
