package internal

import (
	"fmt"
	"time"

	"github.com/Eventid3/peon/database"
)

const TICKER_FREQ = 5 * time.Second

type Scheduler struct {
	Ticker   *time.Ticker
	DBCtx    *database.DatabaseContext
	JobChan  chan WorkerJob
	Registry JobRegistry
}

func NewScheduler(dbCtx *database.DatabaseContext, jobChan chan WorkerJob, registry map[string]JobRegistration) *Scheduler {
	return &Scheduler{
		time.NewTicker(TICKER_FREQ),
		dbCtx,
		jobChan,
		registry,
	}
}

func (s *Scheduler) StartScheduler() {
	go func() {
		for {
			<-s.Ticker.C
			_ = s.EnqueueJobs()
			_ = s.ScheduleJobs()
		}
	}()
}

// add job definitions in "queued"-state to job runs
func (s *Scheduler) EnqueueJobs() error {
	jobs, err := s.DBCtx.GetJobDefinitionsByNextRun(time.Now())
	if err != nil {
		return fmt.Errorf("error enqueing jobs: %w", err)
	}

	errs := make([]error, len(jobs))
	for _, job := range jobs {
		err = s.DBCtx.CreateJobDefinition(job)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("error(s) creating job runs while enqueing: %s", errorsToString(errs))
	}
	return nil
}

// hand jobs to workers
func (s *Scheduler) ScheduleJobs() error {
	return nil
}

func errorsToString(errors []error) string {
	res := ""
	for _, err := range errors {
		res = fmt.Sprintf("    %s\n    %s", res, err)
	}
	return res
}
