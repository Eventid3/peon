// Package internal holds the internal mechanisms of Peon,
// that are not relevant for users of the library
package internal

import (
	"context"

	"github.com/Eventid3/peon/job"
)

type WorkerJob struct {
	JobName string
	Job     job.Job
}

type WorkerResult struct {
	JobName string
	Err     error
}

type Worker struct {
	jobsChan   *chan WorkerJob
	resultChan *chan WorkerResult
}

func NewWorker(jc *chan WorkerJob, rc *chan WorkerResult) Worker {
	return Worker{jobsChan: jc, resultChan: rc}
}

func (w *Worker) DoWork(ctx context.Context) {
	for work := range *w.jobsChan {
		err := work.Job.Exec(ctx)
		*w.resultChan <- WorkerResult{work.JobName, err}
	}
}
