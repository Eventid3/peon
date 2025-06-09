package peon

import (
	"fmt"
	"sync"
	"time"
)

type Workforce struct {
	name    string
	workers map[string]*Worker
	mu      sync.Mutex
}

func NewWorkforce(name string) *Workforce {
	return &Workforce{
		name:    name,
		workers: make(map[string]*Worker),
	}
}

func (wf *Workforce) AddRepeatableJob(name string, interval time.Duration, job Job) {
	_, ok := wf.workers[name]
	if ok {
		fmt.Printf("Worker with name %s already exists\n", name)
		return
	}

	worker := NewWorker(name)
	worker.AddRepeatableJob(name, interval, job)

	wf.workers[name] = worker
}

func (wf *Workforce) StartAllWorkers() {
	for _, worker := range wf.workers {
		worker.StartWorker()
	}
}
