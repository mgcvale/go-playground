package config

import (
	"awesomeProject/models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabse() {
	state := GetApplicationState()

	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	state.SetDB(db)
}

func Migrate() error {
	state := GetApplicationState()
	db := state.GetDB()

	if db == nil {
		return fmt.Errorf("Database is null in ApplicationState. Call InitDatabase() first.")
	}

	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	return nil
}
