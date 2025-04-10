package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ritikchawla/url-shortner/config"
)

var DB *sql.DB

func InitPostgres(cfg *config.Config) error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	// Create URLs table if it doesn't exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			long_url TEXT NOT NULL,
			short_code VARCHAR(10) UNIQUE NOT NULL,
			visits BIGINT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
	`

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}

func ClosePostgres() {
	if DB != nil {
		DB.Close()
	}
}
