package daemon

import (
	"fmt"
	"time"

	"github.com/kycklingar/pultimad/config"
)

type Checker interface {
	Init(config.Config) error
	Check() []Taskif
	Quit()
}

type domain struct {
	checker Checker
	name    string
	queue   *queue

	sleepTime time.Duration

	q bool
}

func (dom *domain) quit() {
	dom.q = true
	dom.checker.Quit()
}

func (dom *domain) process() {
	for !dom.q {
		task := dom.queue.pop()
		if task != nil {
			fmt.Printf("%s - [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), task.Domain(), task.Description())
			newTasks := task.Do()
			if newTasks != nil {
				dom.queue.push(newTasks)
			}
		} else if !dom.check() {
			time.Sleep(time.Second * 10)
		}

		time.Sleep(dom.sleepTime)
	}
}

func (dom *domain) check() bool {
	tasks := dom.checker.Check()
	if tasks != nil {
		dom.queue.push(tasks)
	}

	return len(tasks) > 0
}
