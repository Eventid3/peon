package database

import (
	"context"
	"fmt"
	"uuid"

	"github.com/Eventid3/peon/database/models"
	"github.com/jackc/pgx/v5"
)

func (dbCtx *DatabaseContext) CreateJobDefinition(jd models.JobDefinition) error {
	sql := fmt.Sprintf(`
		INSERT INTO %s (name, timing, retry_count, active, next_run_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (name) DO UPDATE
		SET timing = EXCLUDED.timing,
		retry_count = EXCLUDED.retry_count,
		next_run_at = CASE
			WHEN job_definitions.timing IS DISTINCT FROM EXCLUDED.timing
			THEN EXCLUDED.next_run_at
			ELSE job_definitions.next_run_at
		END;
	`, JOB_DEFINITIONS_TABLE)

	_, err := dbCtx.pool.Exec(context.Background(), sql,
		jd.Name,
		jd.Timing,
		jd.RetryCount,
		jd.Active,
		jd.NextRunAt,
	)
	if err != nil {
		return fmt.Errorf("error inserting job definition: %w", err)
	}

	return nil
}

func (dbCtx *DatabaseContext) UpdateJobDefinition(jd models.JobDefinition) error {
	sql := fmt.Sprintf(`
	UPDATE %s 
	SET timing = $1, retry_count = $2, active = $3
	WHERE name = $4
	`, JOB_DEFINITIONS_TABLE)

	_, err := dbCtx.pool.Exec(context.Background(), sql,
		jd.Timing,
		jd.RetryCount,
		jd.Active,
		jd.Name,
	)
	if err != nil {
		return fmt.Errorf("error updating job definition: %w", err)
	}

	return nil
}

func (dbCtx *DatabaseContext) GetJobDefinitionByName(name string) (models.JobDefinition, error) {
	sql := fmt.Sprintf(`
	SELECT id, name, timing, retry_count, active FROM %s
	WHERE name = $1
	`, JOB_DEFINITIONS_TABLE)

	rows, err := dbCtx.pool.Query(context.Background(), sql,
		name,
	)
	if err != nil {
		return models.JobDefinition{}, fmt.Errorf("error querying job definition: %w", err)
	}

	jobDef, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.JobDefinition])
	if err != nil {
		return models.JobDefinition{}, fmt.Errorf("error collecting job definition: %w", err)
	}
	return jobDef, nil
}

func (dbCtx *DatabaseContext) GetJobDefinitionByID(id uuid.UUID) (models.JobDefinition, error) {
	sql := fmt.Sprintf(`
	SELECT id, name, timing, retry_count, active, next_run_at FROM %s
	WHERE id = $2
	`, JOB_DEFINITIONS_TABLE)

	rows, err := dbCtx.pool.Query(context.Background(), sql,
		id,
	)
	if err != nil {
		return models.JobDefinition{}, fmt.Errorf("error querying job definition: %w", err)
	}

	jobDef, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.JobDefinition])
	if err != nil {
		return models.JobDefinition{}, fmt.Errorf("error collecting job definition: %w", err)
	}
	return jobDef, nil
}

func (dbCtx *DatabaseContext) GetAllJobDefinitions() ([]models.JobDefinition, error) {
	sql := fmt.Sprintf(`
	SELECT id, name, timing, retry_count, next_run_at, active
	FROM %s
	`, JOB_DEFINITIONS_TABLE)

	rows, err := dbCtx.pool.Query(context.Background(), sql)
	if err != nil {
		return []models.JobDefinition{}, fmt.Errorf("error querying job definitions table: %w", err)
	}

	jobDefinitions, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.JobDefinition])
	if err != nil {
		return []models.JobDefinition{}, fmt.Errorf("error collecting job definitions rows: %w", err)
	}

	return jobDefinitions, nil
}

func (dbCtx *DatabaseContext) DeleteJobDefinition(jd models.JobDefinition) error {
	return dbCtx.DeleteJobDefinitionByName(jd.Name)
}

func (dbCtx *DatabaseContext) DeleteJobDefinitionByID(id uuid.UUID) error {
	sql := fmt.Sprintf(`
	DELETE FROM %s WHERE id = $1
	`, JOB_DEFINITIONS_TABLE)

	_, err := dbCtx.pool.Exec(context.Background(), sql,
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting job definition: %w", err)
	}
	return nil
}

func (dbCtx *DatabaseContext) DeleteJobDefinitionByName(name string) error {
	sql := fmt.Sprintf(`
	DELETE FROM %s WHERE name = $1
	`, JOB_DEFINITIONS_TABLE)

	_, err := dbCtx.pool.Exec(context.Background(), sql,
		name,
	)
	if err != nil {
		return fmt.Errorf("error deleting job definition: %w", err)
	}
	return nil
}
