package database

import (
	"context"
	"fmt"
	"uuid"

	"github.com/Eventid3/peon/database/models"
)

func (dbCtx *DatabaseContext) CreateJobDefinition(jd models.JobDefinition) error {
	conn, err := dbCtx.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	sql := fmt.Sprintf(`
		INSERT INTO %s (name, timing, retry_count, active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO UPDATE
		SET timing = EXCLUDED.timing,
		retry_count = EXCLUDED.retry_count
	`, JOB_DEFINITIONS_TABLE)

	_, err = conn.Exec(context.Background(), sql,
		jd.Name,
		jd.Timing,
		jd.RetryCount,
		jd.Active,
	)
	if err != nil {
		return fmt.Errorf("error inserting job definition: %w", err)
	}

	return nil
}

func (dbCtx *DatabaseContext) UpdateJobDefinition(jd models.JobDefinition) error {
	conn, err := dbCtx.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	sql := `
	UPDATE $1 
	SET timing = $2, retry_count = $3, active = $4
	WHERE name = $5
	`

	_, err = conn.Exec(context.Background(), sql,
		JOB_DEFINITIONS_TABLE,
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

func (dbCtx *DatabaseContext) GetJobDefinitionByName(name string) error {
	conn, err := dbCtx.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	sql := `
	SELECT id, name, timing, retry_count, active
	FROM $1 WHERE name = $2
	`

	var jobDef models.JobDefinition

	err = conn.QueryRow(context.Background(), sql,
		JOB_DEFINITIONS_TABLE,
		name,
	).Scan(&jobDef)
	if err != nil {
		return fmt.Errorf("error updating job definition: %w", err)
	}
	return nil
}

func (dbCtx *DatabaseContext) GetJobDefinitionByID(id uuid.UUID) error {
	conn, err := dbCtx.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	sql := `
	SELECT id, name, timing, retry_count, active
	FROM $1 WHERE id = $2
	`

	var jobDef models.JobDefinition

	err = conn.QueryRow(context.Background(), sql,
		JOB_DEFINITIONS_TABLE,
		id,
	).Scan(&jobDef)
	if err != nil {
		return fmt.Errorf("error updating job definition: %w", err)
	}
	return nil
}

func (dbCtx *DatabaseContext) DeleteJobDefinition(jd models.JobDefinition) error {
	return dbCtx.DeleteJobDefinitionByName(jd.Name)
}

func (dbCtx *DatabaseContext) DeleteJobDefinitionByID(id uuid.UUID) error {
	conn, err := dbCtx.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	sql := `
	DELETE FROM $1 WHERE id = $2
	`

	_, err = conn.Exec(context.Background(), sql,
		JOB_DEFINITIONS_TABLE,
		id,
	)
	if err != nil {
		return fmt.Errorf("error deleting job definition: %w", err)
	}
	return nil
}

func (dbCtx *DatabaseContext) DeleteJobDefinitionByName(name string) error {
	conn, err := dbCtx.Connect()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	sql := `
	DELETE FROM $1 WHERE name = $2
	`

	_, err = conn.Exec(context.Background(), sql,
		JOB_DEFINITIONS_TABLE,
		name,
	)
	if err != nil {
		return fmt.Errorf("error deleting job definition: %w", err)
	}
	return nil
}
