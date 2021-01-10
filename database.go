package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectDatabase() (*gorm.DB, error) {
	database, issue := gorm.Open(sqlite.Open("database.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if issue != nil {
		return nil, issue
	}

	database.AutoMigrate(&Account{})
	return database, nil
}
