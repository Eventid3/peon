// Package models holds the database models of the Peon
// job scheduler
package models

import "uuid"

type JobDefinition struct {
	ID         uuid.UUID
	Name       string
	Timing     string
	RetryCount int
	Active     bool
}
