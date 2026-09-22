package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const (
	TABLE_NAME            = "peon_db"
	JOB_DEFINITIONS_TABLE = "job_definitions"
	JOB_RUNS_TABLE        = "job_runs"
)

type DatabaseContext struct {
	baseConnectionString string
	connectionString     string
}

func NewDatabaseContext(baseConnStr string) *DatabaseContext {
	connStr := fmt.Sprintf("%s/%s?sslmode=disable", baseConnStr, TABLE_NAME)
	dbCtx := DatabaseContext{baseConnectionString: baseConnStr, connectionString: connStr}
	err := dbCtx.initDatabase()
	if err != nil {
		fmt.Println("error initializing database: %w", err)
	}
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

	err = conn.QueryRow(context.Background(), sql, TABLE_NAME).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error initializing the database: %w", err)
	}

	if !exists {
		fmt.Println("[Database] creating database...")
		_, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", TABLE_NAME))
		if err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}
	} else {
		fmt.Println("[Database] database already exists...")
	}

	return nil
}

func (dbCtx *DatabaseContext) initTables() error {
	fmt.Println("[Database] initializing tables...")
	conn, err := dbCtx.Connect()
	if err != nil {
		return fmt.Errorf("error connecting to the database context: %w", err)
	}
	defer func() {
		_ = conn.Close(context.Background())
	}()

	sql := `SELECT EXISTS(
SELECT 1 FROM pg_catalog.pg_tables
WHERE schemaname != 'pg_catalog' AND
schemaname != 'information_schema' AND
tablename = $1)`

	var jobDefinitionsExists bool
	err = conn.QueryRow(context.Background(), sql, JOB_DEFINITIONS_TABLE).Scan(&jobDefinitionsExists)
	if err != nil {
		return fmt.Errorf("error scanning for job definitions table: %w", err)
	}

	if !jobDefinitionsExists {
		fmt.Println("[Database] creating job definitions table..")
		_, err = conn.Exec(context.Background(),
			fmt.Sprintf(`
				CREATE TABLE %s (
				id uuid primary key,
				name text,
				timing varchar(16),
				retry_count int,
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
	err = conn.QueryRow(context.Background(), sql, JOB_RUNS_TABLE).Scan(&jobRunsExists)
	if err != nil {
		return fmt.Errorf("error scanning for job runs table: %w", err)
	}
	// TODO create table
	if !jobRunsExists {
		fmt.Println("[Database] creating job runs table...")
		_, err = conn.Exec(context.Background(),
			fmt.Sprintf(`
				CREATE TABLE %s (
				id uuid primary key,
				job_id uuid references %s(id),
				status text,
				next_retry_at text,
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

func (dbCtx *DatabaseContext) Connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(
		context.Background(),
		dbCtx.connectionString,
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
