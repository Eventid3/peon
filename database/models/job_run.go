package models

import (
	"time"
	"uuid"
)

type JobRun struct {
	ID           uint32
	JobID        uuid.UUID
	Status       string
	NextRetryAt  time.Time
	AttemptCount int
	Error        string
}
