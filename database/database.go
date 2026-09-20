package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DatabaseContext struct {
	connectionString string
}

func NewDatabaseContext(connStr string) *DatabaseContext {
	dbCtx := DatabaseContext{connectionString: connStr}
	err := dbCtx.init()
	if err != nil {
		fmt.Println("error initializing database: %w", err)
	}
	return &dbCtx
}

func (dbCtx *DatabaseContext) init() error {
	conn, err := dbCtx.Connect()
	if err != nil {
		return fmt.Errorf("error connecting to the database context: %w", err)
	}
	defer conn.Close(context.Background())

	sql := `
		SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)
	`

	var exists bool

	err = conn.QueryRow(context.Background(), sql, "peon_db").Scan(&exists)
	if err != nil {
		return fmt.Errorf("error initializing the database: %w", err)
	}

	if !exists {
		_, err = conn.Exec(context.Background(), "CREATE DATABASE peon_db")
		if err != nil {
			return fmt.Errorf("error creating database: %w", err)
		}
	}

	if err != nil {
		return fmt.Errorf("error closing database connection: %w", err)
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
