package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// InitDB initializes the database connection.
//
// It takes a database path as a parameter and returns a *gorm.DB and an error.
func InitDB(databasePath string) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = DB.AutoMigrate(&Guild{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate tables: %w", err)
	}

	return DB, nil
}
