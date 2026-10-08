// Package models holds the database models of the Peon
// job scheduler
package models

import (
	"database/sql"
	"time"
	"uuid"
)

type JobDefinition struct {
	ID         uuid.UUID
	Name       string
	Timing     time.Duration
	NextRunAt  sql.NullTime
	RetryCount int
	Active     bool
}
