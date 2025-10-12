package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connects to database. If db file does not exist it creates a new recipes.db file
func ConnectToSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("recipes.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
