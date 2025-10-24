package database

import (
	"database/sql/driver"
	"encoding/json"
	"recipeapp/client"
	"recipeapp/models"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MealsJSON []models.Meal

type RecipesEntry struct {
	EntryUUID uuid.UUID `gorm:"primaryKey"`
	Meals     MealsJSON `gorm:"type:json"`
}

// Value marshals the MealsJSON slice into a JSON byte array for database storage
func (m MealsJSON) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan unmarshals JSON data from the database back into a MealsJSON slice
func (m *MealsJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

// ConnectToSQLite opens (or creates) the SQLite database and returns the DB instance
func ConnectToSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("recipes.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateResult creates a new RecipesEntry in the database and returns its UUID
func CreateResult(db *gorm.DB, response *client.Response) (uuid.UUID, error) {
	entryUUID := uuid.New()
	meals := response.Meals
	entry := RecipesEntry{
		EntryUUID: entryUUID,
		Meals:     meals,
	}
	if err := db.Create(&entry).Error; err != nil {
		return uuid.Nil, err
	}
	return entryUUID, nil
}

// GetRecipesFromDBByUUID retrieves a RecipesEntry by UUID and returns its Meals slice
func GetRecipesFromDBByUUID(db *gorm.DB, id uuid.UUID) ([]models.Meal, error) {
	var entry RecipesEntry
	if err := db.First(&entry, "entry_uuid = ?", id).Error; err != nil {
		return nil, err
	}
	return entry.Meals, nil
}
