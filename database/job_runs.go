package database

import (
	"context"
	"fmt"
	"uuid"

	"github.com/Eventid3/peon/database/models"
)

func (dbCtx *DatabaseContext) CreateJobRun(jd models.JobDefinition) error {
	sql := fmt.Sprintf(`
		INSERT INTO %s (job_id, status, attempt_count)
		VALUES ($1, $2, $3)
		`, JOB_RUNS_TABLE)

	_, err := dbCtx.pool.Exec(context.Background(), sql, jd.ID, models.Queued.String(), 0)
	if err != nil {
		return fmt.Errorf("error queuing job in database: %w", err)
	}
	return nil
}

func (dbCtx *DatabaseContext) GetJobRunsByName(jobName string) ([]models.JobRun, error) {
	return []models.JobRun{}, nil
}

func (dbCtx *DatabaseContext) GetJobsToRun() ([]uuid.UUID, error) {
	return []uuid.UUID{}, nil
}

func (dbCtx *DatabaseContext) UpdateJobRun(jr models.JobRun) error {
	return nil
}

func (dbCtx *DatabaseContext) DeleteJobRunByID(id uint32) error {
	return nil
}

func (dbCtx *DatabaseContext) ClearJobRuns() error {
	return nil
}
