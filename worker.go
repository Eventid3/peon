package peon

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

type Job interface {
	Execute(context.Context) error
}

type WorkerJob struct {
	Name     string
	Interval time.Duration
	Timeout  time.Duration
	Active   bool
	Job      Job
}

type Worker struct {
	jobs   map[string]*WorkerJob
	logger *log.Logger
	ctx    context.Context
	cancel context.CancelFunc
	mutex  sync.Mutex
	wg     sync.WaitGroup
}

func NewWorker() *Worker {
	return &Worker{
		jobs:   make(map[string]*WorkerJob),
		logger: log.New(os.Stdout, "[Worker] ", log.Ltime),
	}
}

func (w *Worker) AddRepeatableJob(name string, interval time.Duration, timeout time.Duration, job Job) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	_, j := w.jobs[name]
	if j {
		w.logger.Printf("Cannot add job with name %v: name already exists.", name)
		return
	}

	workerJob := WorkerJob{
		Name:     name,
		Interval: interval,
		Timeout:  timeout,
		Active:   false,
		Job:      job,
	}

	w.jobs[name] = &workerJob
}

func (w *Worker) StartWorker(ctx context.Context) {
	w.ctx, w.cancel = context.WithCancel(ctx)

	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, wj := range w.jobs {
		wj.Active = true
		w.wg.Add(1)
		go w.runJob(wj)

	}
}

func (w *Worker) runJob(wj *WorkerJob) {
	defer w.wg.Done()

	ticker := time.NewTicker(wj.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():

			w.logger.Printf("Worker %v stopped, due to timeout\n", wj.Name)
		case <-ticker.C:
			if !wj.Active {
				continue
			}
			w.executeJob(wj)
		}
	}
}

func (w *Worker) executeJob(wj *WorkerJob) {
	w.logger.Printf("Executing job %v\n", wj.Name)

	jobCtx := w.ctx
	var cancel context.CancelFunc

	if wj.Timeout > 0 {
		jobCtx, cancel = context.WithTimeout(jobCtx, wj.Timeout)
		defer cancel()
	}

	done := make(chan error, 1)
	go func() {
		done <- wj.Job.Execute(jobCtx)
	}()

	select {
	case err := <-done:
		if err != nil {
			w.logger.Printf("Job %v failed: %v\n", wj.Name, err)
		} else {
			w.logger.Printf("Job %v finished successfully\n", wj.Name)
		}
	case <-jobCtx.Done():
		w.logger.Printf("Job %v timed out\n", wj.Name)
	}
}

func (w *Worker) StopWorker() {
	w.logger.Printf("Stopping worker\n")
	if w.cancel != nil {
		w.cancel()
	}
	w.wg.Wait()
	w.logger.Printf("Worker stopped\n")
}

func (w *Worker) StopJob(name string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if job, exists := w.jobs[name]; exists {
		job.Active = false
		w.logger.Printf("Job %v stopped\n", name)
	}
}

func (w *Worker) StartJob(name string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if job, exists := w.jobs[name]; exists {
		job.Active = true
		w.logger.Printf("Job %v started\n", name)
	}
}
