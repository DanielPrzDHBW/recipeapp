package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectToSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("recipes.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
