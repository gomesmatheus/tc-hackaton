package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

const (
	createTables = `
		CREATE TABLE IF NOT EXISTS videos (
			id VARCHAR(255) PRIMARY KEY,
			owner_id VARCHAR(255) NOT NULL,
			status VARCHAR(20) NOT NULL
		);
    `
)

func NewPostgresDb(url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		fmt.Println("Error parsing config", err)
	}
	db, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Println("Error creating database connection", err)
		return nil, err
	}

	if _, err := db.Exec(context.Background(), createTables); err != nil {
		fmt.Println("Error creating table Persons", err)
		return nil, err
	}

	return db, err
}
