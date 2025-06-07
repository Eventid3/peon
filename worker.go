package peon

import (
	"context"
	"log"
	"os"
	"time"
)

type WorkerJob struct {
	Name     string
	Interval time.Duration
	Active   bool
	Job      Job
}

type Worker struct {
	workerJob *WorkerJob
	logger    *log.Logger
	ctx       context.Context
	cancel    context.CancelFunc
	ticker    *time.Ticker
	// wg        sync.WaitGroup
}

func NewWorker() *Worker {
	return &Worker{
		logger: log.New(os.Stdout, "[Worker] ", log.Ltime),
	}
}

func (w *Worker) AddRepeatableJob(name string, interval time.Duration, job Job) {
	workerJob := WorkerJob{
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
				w.logger.Printf("Worker %v paused\n", w.workerJob.Name)
				continue
			}
			w.logger.Printf("Starting job %v\n", w.workerJob.Name)
			err := w.workerJob.Job.Execute()
			if err != nil {
				w.logger.Printf("Error while running job %v: %v\n", w.workerJob.Name, err)
			} else {
				w.logger.Printf("Job %v completed successfully\n", w.workerJob.Name)
			}
		case <-w.ctx.Done():
			w.logger.Printf("Worker %v stopped, due to context cancellation\n", w.workerJob.Name)
			w.workerJob.Active = false
			return
		}
	}
}

func (w *Worker) PauseWorker() {
	w.logger.Printf("Pausing worker %v\n", w.workerJob.Name)
	w.workerJob.Active = false
	w.ticker.Stop()
}

func (w *Worker) ResumeWorker() {
	w.logger.Printf("Resuming worker %v\n", w.workerJob.Name)
	w.workerJob.Active = true
	w.ticker.Reset(w.workerJob.Interval)
}

func (w *Worker) StopWorker() {
	w.logger.Printf("Stopping worker\n")
	if w.cancel != nil {
		w.cancel()
	}
	w.logger.Printf("Worker stopped\n")
}
