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
	// Make sure the directory db exists
	if _, err := os.Stat("db"); os.IsNotExist(err) {
		os.Mkdir("db", 0755)
	}
	if file == "" {
		file = "db/dev.sqlite"
	}
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &service{
		db: db,
	}
}
