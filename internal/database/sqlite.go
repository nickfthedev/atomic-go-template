package database

import (
	"log"
	"os"

	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

var (
	//SQLite
	file = os.Getenv("DB_FILE")
)

func NewSQLiteService() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_FILE")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &service{
		db: db,
	}
}
