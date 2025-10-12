package database

import (
	"github.com/google/uuid"
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

type RecipesEntry struct {
	UUID     string
	DataJSON string
}

// Writes a new Result to DB under UUID and returns the uuid, so that it can be stored for access to the data
func CreateResult(db *gorm.DB, dataJSON string) (string, error) {
	uuid := generateUUID()
	entry := RecipesEntry{
		uuid,
		dataJSON,
	}
	db.Create(&entry)
	return uuid, nil
}

// Generates new UUID for new DB entry
func generateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}
