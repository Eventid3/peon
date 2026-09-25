package database

import (
	"context"
	"fmt"
	"uuid"

	"github.com/Eventid3/peon/database/models"
)

type JobState int

const (
	Queued JobState = iota
	Running
	Retrying
	Failed
	Succeeded
)

func (js JobState) String() string {
	switch js {
	case 0:
		return "Queued"
	case 1:
		return "Running"
	case 2:
		return "Retrying"
	case 3:
		return "Failed"
	case 4:
		return "Succeeded"
	}
	return fmt.Sprintf("JobState(%v)", int(js))
}

func (dbCtx *DatabaseContext) CreateJobRun(jd models.JobDefinition) error {
	conn, err := dbCtx.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	sql := fmt.Sprintf(`
		INSERT INTO %s (job_id, status, attempt_count)
		VALUES ($1, $2, $3)
		`, JOB_RUNS_TABLE)

	_, err = conn.Exec(context.Background(), sql, jd.ID, Queued.String(), 0)
	if err != nil {
		return fmt.Errorf("error queuing job in database: %w", err)
	}
	return nil
}

func (dbCtx *DatabaseContext) GetJobRunsByName(jobName string) ([]models.JobRun, error) {
	conn, err := dbCtx.Connect()
	if err != nil {
		return []models.JobRun{}, err
	}
	defer conn.Close(context.Background())
	return []models.JobRun{}, nil
}

func (dbCtx *DatabaseContext) GetJobsToRun() ([]uuid.UUID, error) {
	conn, err := dbCtx.Connect()
	if err != nil {
		return []uuid.UUID{}, err
	}
	defer conn.Close(context.Background())
	return []uuid.UUID{}, nil
}

func (dbCtx *DatabaseContext) UpdateJobRun(jr models.JobRun) error {
	conn, err := dbCtx.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	return nil
}

func (dbCtx *DatabaseContext) DeleteJobRunByID(id uint32) error {
	conn, err := dbCtx.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	return nil
}

func (dbCtx *DatabaseContext) ClearJobRuns() error {
	conn, err := dbCtx.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	return nil
}
