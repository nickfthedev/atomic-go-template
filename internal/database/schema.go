package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MigrateUserSchema migrates the user schema to the database.
func MigrateUserSchema(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

// BaseModel includes common fields for all tables.
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}

// User represents a user in the database.
type User struct {
	BaseModel
	Username      string  `gorm:"unique;not null"`
	Email         string  `gorm:"unique;not null"`
	Password      *string `gorm:""` // Password is optional
	OAuthProvider *string `gorm:""` // OAuth provider name (e.g., "google", "github")
	OAuthID       *string `gorm:""` // OAuth provider user ID
}
