// Package database
package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	DATABASE_NAME         = "peon_db"
	JOB_DEFINITIONS_TABLE = "job_definitions"
	JOB_RUNS_TABLE        = "job_runs"
)

type DatabaseContext struct {
	baseConnectionString string
	connectionString     string
	pool                 *pgxpool.Pool
}

func NewDatabaseContext(baseConnStr string) *DatabaseContext {
	connStr := fmt.Sprintf("%s/%s?sslmode=disable", baseConnStr, DATABASE_NAME)
	dbCtx := DatabaseContext{baseConnectionString: baseConnStr, connectionString: connStr}
	err := dbCtx.initDatabase()
	if err != nil {
		fmt.Println("error initializing database: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Println("error creating connection pool: %w", err)
	}
	dbCtx.pool = pool

	err = dbCtx.initTables()
	if err != nil {
		fmt.Println("error initializing tables: %w", err)
	}
	return &dbCtx
}

func (dbCtx *DatabaseContext) initDatabase() error {
	fmt.Println("[Database] initializing databse...")
	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf("%s/postgres?sslmode=disable", dbCtx.baseConnectionString),
	)
	if err != nil {
		return fmt.Errorf("error connecting to the database context: %w", err)
	}
	defer func() {
		_ = conn.Close(context.Background())
	}()

	sql := `
		SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)
	`

	var exists bool

	err = conn.QueryRow(context.Background(), sql, DATABASE_NAME).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error initializing the database: %w", err)
	}

	if !exists {
		fmt.Println("[Database] creating database...")
		_, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", DATABASE_NAME))
		if err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}
	} else {
		fmt.Println("[Database] database already exists...")
	}

	return nil
}

func (dbCtx *DatabaseContext) initTables() error {
	sql := `SELECT EXISTS(
	SELECT 1 FROM pg_catalog.pg_tables
	WHERE schemaname != 'pg_catalog' AND
	schemaname != 'information_schema' AND
	tablename = $1)`

	var jobDefinitionsExists bool
	err := dbCtx.pool.QueryRow(context.Background(), sql, JOB_DEFINITIONS_TABLE).Scan(&jobDefinitionsExists)
	if err != nil {
		return fmt.Errorf("error scanning for job definitions table: %w", err)
	}

	if !jobDefinitionsExists {
		fmt.Println("[Database] creating job definitions table..")
		_, err = dbCtx.pool.Exec(context.Background(),
			fmt.Sprintf(`
				CREATE TABLE %s (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				name text UNIQUE,
				timing interval,
				next_run_at timestamp,
				retry_count int CHECK (retry_count > 0),
				active boolean
				)
				`, JOB_DEFINITIONS_TABLE),
		)
		if err != nil {
			return fmt.Errorf("error creating job definitions table: %w", err)
		}
	} else {
		fmt.Println("[Database] job definitions table already exists...")
	}

	var jobRunsExists bool
	err = dbCtx.pool.QueryRow(context.Background(), sql, JOB_RUNS_TABLE).Scan(&jobRunsExists)
	if err != nil {
		return fmt.Errorf("error scanning for job runs table: %w", err)
	}

	if !jobRunsExists {
		fmt.Println("[Database] creating job runs table...")
		_, err = dbCtx.pool.Exec(context.Background(),
			fmt.Sprintf(`
				CREATE TABLE %s (
				id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
				job_id uuid REFERENCES %s(id),
				status text NOT NULL,
				started_at timestamp,
				next_retry_at timestamp,
				completed_at timestamp,
				attempt_count int,
				error text
				)
				`, JOB_RUNS_TABLE, JOB_DEFINITIONS_TABLE),
		)
		if err != nil {
			return fmt.Errorf("error creating job definitions table: %w", err)
		}
	} else {
		fmt.Println("[Database] job table already exists...")
	}

	return nil
}
