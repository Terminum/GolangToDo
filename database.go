package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func connectDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=admin dbname=todo sslmode=disable password=pass"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}
