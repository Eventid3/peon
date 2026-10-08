package internal

import (
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
	return nil
}

// hand jobs to workers
func (s *Scheduler) ScheduleJobs() error {
	return nil
}
