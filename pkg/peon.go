// Package pkg holds the main entrypoint for starting and configuring
// a new Peon instance and registering jobs
package pkg

import (
	"container/list"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Eventid3/peon/database"
	"github.com/Eventid3/peon/database/models"
	"github.com/Eventid3/peon/internal"
	"github.com/Eventid3/peon/job"
)

type Peon struct {
	workerPool   []internal.Worker
	queue        *list.List
	databaseCtx  *database.DatabaseContext
	jobRegistry  internal.JobRegistry
	jobChan      chan internal.WorkerJob
	responseChan chan internal.WorkerResult
	scheduler    *internal.Scheduler
	running      bool
}

func NewPeon(numWorkers int) *Peon {
	databaseCtx := database.NewDatabaseContext("postgres://admin:password@localhost:5432")
	jobChan := make(chan internal.WorkerJob)
	responseChan := make(chan internal.WorkerResult)
	workerPool := make([]internal.Worker, numWorkers)
	jobRegistry := internal.NewJobRegistry()

	for i := range workerPool {
		workerPool[i] = internal.NewWorker(&jobChan, &responseChan)
	}

	return &Peon{
		workerPool:   workerPool,
		queue:        list.New(),
		databaseCtx:  databaseCtx,
		jobRegistry:  jobRegistry,
		jobChan:      jobChan,
		responseChan: responseChan,
		scheduler:    internal.NewScheduler(databaseCtx, jobChan, jobRegistry),
		running:      false,
	}
}

func (p *Peon) RegisterJob(jobName string, job job.Job, timing time.Duration) error {
	_, ok := p.jobRegistry[jobName]
	if ok {
		return fmt.Errorf("could not register %s - job already exists in job registry", jobName)
	}

	err := p.databaseCtx.CreateJobDefinition(models.JobDefinition{
		Name:       jobName,
		Timing:     timing,
		Active:     true,
		RetryCount: 1,
		NextRunAt:  sql.NullTime{Time: time.Now().Add(timing), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error registering job %s: %w", jobName, err)
	}

	p.jobRegistry[jobName] = internal.JobRegistration{jobName, job, timing}
	return nil
}

func (p *Peon) Start() {
	p.running = true

	go p.startWorkers()
	go p.scheduleJobs()
	go p.handleResults()
	p.scheduler.StartScheduler()
}

func (p *Peon) Stop() {
	p.running = false
}

func (p *Peon) startWorkers() {
	for _, worker := range p.workerPool {
		go worker.DoWork(context.Background())
	}
}

func (p *Peon) scheduleJobs() {
	ticker := time.NewTicker(3 * time.Second)
	for p.running {
		<-ticker.C
		for name, jobReg := range p.jobRegistry {
			fmt.Printf("starting job %s\n", name)
			p.jobChan <- internal.WorkerJob{JobName: name, Job: jobReg.Job}
		}
	}
}

func (p *Peon) handleResults() {
	for p.running {
		res := <-p.responseChan
		if res.Err != nil {
			fmt.Printf("ERROR: job %s resulted in an error: %s\n", res.JobName, res.Err)
		} else {
			fmt.Printf("SUCCESS: job %s ran successfully\n", res.JobName)
		}
	}
}
