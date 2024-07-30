package database

import (
	"my-go-template/internal/model"

	"gorm.io/gorm"
)

// MigrateUserSchema migrates the user schema to the database.
func MigrateUserSchema(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{})
}

// Models are in the models folder
