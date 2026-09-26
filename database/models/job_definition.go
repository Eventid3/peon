// Package models holds the database models of the Peon
// job scheduler
package models

import (
	"time"
	"uuid"
)

type JobDefinition struct {
	ID         uuid.UUID
	Name       string
	Timing     time.Duration
	NextRunAt  time.Time
	RetryCount int
	Active     bool
}
