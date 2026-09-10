package peon

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

type WorkerJob struct {
	Exec func() error
}

type FireAndForgetJob struct {
	Name string
	Job  WorkerJob
}

type RepeatableWorkerJob struct {
	Name     string
	Interval time.Duration
	Active   bool
	Job      WorkerJob
}

type Worker struct {
	workerJob *RepeatableWorkerJob
	logger    *log.Logger
	ctx       context.Context
	cancel    context.CancelFunc
	ticker    *time.Ticker
	// wg        sync.WaitGroup
}

func NewWorker(name string) *Worker {
	return &Worker{
		logger: log.New(os.Stdout, fmt.Sprintf("[Worker: %s] ", name), log.Ltime),
	}
}

func (w *Worker) AddRepeatableJob(name string, interval time.Duration, job WorkerJob) {
	workerJob := RepeatableWorkerJob{
		Name:     name,
		Interval: interval,
		Active:   false,
		Job:      job,
	}

	w.workerJob = &workerJob
}

func (w *Worker) StartWorker() {
	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.workerJob.Active = true
	w.ticker = time.NewTicker(w.workerJob.Interval)

	go w.runJob()
}

func (w *Worker) runJob() {
	for {
		select {
		case <-w.ticker.C:
			if !w.workerJob.Active {
				w.logger.Println("Worker paused")
				continue
			}
			w.logger.Println("Starting job")
			err := w.workerJob.Job.Exec()
			if err != nil {
				w.logger.Printf("Error while running job: %v\n", err)
			} else {
				w.logger.Println("Job completed successfully")
			}
		case <-w.ctx.Done():
			w.logger.Println("Worker stopped, due to context cancellation")
			w.workerJob.Active = false
			return
		}
	}
}

func (w *Worker) PauseWorker() {
	w.logger.Println("Pausing worker")
	w.workerJob.Active = false
	w.ticker.Stop()
}

func (w *Worker) ResumeWorker() {
	w.logger.Println("Resuming worker")
	w.workerJob.Active = true
	w.ticker.Reset(w.workerJob.Interval)
}

func (w *Worker) StopWorker() {
	w.logger.Println("Stopping worker")
	if w.cancel != nil {
		w.cancel()
	}
	w.logger.Println("Worker stopped")
}
