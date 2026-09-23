package models

import (
	"time"
	"uuid"
)

type JobRun struct {
	ID           uuid.UUID
	JobID        uuid.UUID
	Status       string
	NextRetryAt  time.Time
	AttemptCount int
	Error        string
}
