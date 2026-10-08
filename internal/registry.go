package internal

import (
	"time"

	"github.com/Eventid3/peon/job"
)

type JobRegistration struct {
	Name   string
	Job    job.Job
	Timing time.Duration
}

type JobRegistry = map[string]JobRegistration

func NewJobRegistry() JobRegistry {
	return make(map[string]JobRegistration)
}
