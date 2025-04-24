package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("storage.db"))
	if err != nil {
		log.Fatal("Failed to connect to the DB")
	}
	db.AutoMigrate(&User{},&Book{})
	return db
}
