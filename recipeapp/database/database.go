package database

import (
	"recipeapp/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ConnectToSQLite Connects to database. If db file does not exist it creates a new recipes.db file
func ConnectToSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("recipes.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

type RecipesEntry struct {
	EntryUUID string
	DataJSON  string
}

// TODO: The GetRecipesFromDBByUUID-Function looks for []models.Meal{} in the database.
// This function however stores DataJSON as a string in the database.
// Resolving this issue depends on the datatype that is provided by api.go

// CreateResult Writes a new Result to DB under UUID and returns the uuid, so that it can be stored for access to the data
func CreateResult(db *gorm.DB, dataJSON string) (string, error) {
	entryUUID := generateUUID()
	entry := RecipesEntry{
		entryUUID,
		dataJSON,
	}
	db.Create(&entry)
	if db.Error != nil {
		return "", db.Error
	}
	return entryUUID, nil
}

// Generates new UUID for new DB entry
func generateUUID() string {
	entryUUID := uuid.New()
	return entryUUID.String()
}

// GetRecipesFromDBByUUID Reads the saved data from the DB and returns it
func GetRecipesFromDBByUUID(db *gorm.DB, uuid string) ([]models.Meal, error) {
	result := db.First(&[]models.Meal{}, "uuid = ?", uuid)
	if result.Error != nil {
		return []models.Meal{}, result.Error
	}
	return []models.Meal{}, nil
}
