package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// ConnectDB connects to the SQLite DB
func ConnectDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	return db, nil
}
