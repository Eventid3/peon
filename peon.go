package peon

import (
	"container/list"
	"context"
	"fmt"
	"time"

	"github.com/Eventid3/peon/internal"
	"github.com/Eventid3/peon/job"
)

type jobRegistration struct {
	name   string
	job    job.Job
	timing time.Duration
}

type Peon struct {
	workerPool   []internal.Worker
	queue        *list.List
	jobRegistry  map[string]jobRegistration
	jobChan      chan internal.WorkerJob
	responseChan chan internal.WorkerResult
	running      bool
}

func NewPeon(numWorkers int) *Peon {
	jobChan := make(chan internal.WorkerJob)
	responseChan := make(chan internal.WorkerResult)
	workerPool := make([]internal.Worker, numWorkers)

	for i := range workerPool {
		workerPool[i] = internal.NewWorker(&jobChan, &responseChan)
	}

	return &Peon{
		workerPool:   workerPool,
		queue:        list.New(),
		jobRegistry:  make(map[string]jobRegistration),
		jobChan:      jobChan,
		responseChan: responseChan,
		running:      false,
	}
}

func (p *Peon) RegisterJob(jobName string, job job.Job, timing time.Duration) error {
	_, ok := p.jobRegistry[jobName]
	if ok {
		return fmt.Errorf("could not register %s - job already exists in job registry", jobName)
	}
	p.jobRegistry[jobName] = jobRegistration{jobName, job, timing}
	return nil
}

func (p *Peon) Start() {
	p.running = true

	go p.startWorkers()
	go p.scheduleJobs()
	go p.handleResults()
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
			p.jobChan <- internal.WorkerJob{JobName: name, Job: jobReg.job}
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
