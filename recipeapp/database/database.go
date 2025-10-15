package database

import (
	"database/sql/driver"
	"encoding/json"
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

func (m MealsJSON) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MealsJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, m)
}

// ConnectToSQLite Connects to database. If db file does not exist it creates a new recipes.db file
func ConnectToSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("recipes.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateResult Writes a new Result to DB under UUID and returns the uuid, so that it can be stored for access to the data
func CreateResult(db *gorm.DB, meals []models.Meal) (uuid.UUID, error) {
	entryUUID := uuid.New()
	entry := RecipesEntry{
		EntryUUID: entryUUID,
		Meals:     meals,
	}
	if err := db.Create(&entry).Error; err != nil {
		return uuid.Nil, err
	}
	return entryUUID, nil
}

// GetRecipesFromDBByUUID Reads the saved data from the DB and returns it
func GetRecipesFromDBByUUID(db *gorm.DB, id uuid.UUID) ([]models.Meal, error) {
	var entry RecipesEntry
	if err := db.First(&entry, "entry_uuid = ?", id).Error; err != nil {
		return nil, err
	}
	return entry.Meals, nil
}
