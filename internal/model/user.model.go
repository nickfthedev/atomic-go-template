package model

import "time"

// User represents a user in the database.
type User struct {
	BaseModel
	Username      string     `gorm:"unique;not null"`
	Email         string     `gorm:"unique;not null"`
	Password      *string    `gorm:""` // Password is optional
	VerifiedAt    *time.Time `gorm:""` // Verified at is optional
	OAuthProvider *string    `gorm:""` // OAuth provider name (e.g., "google", "github")
	OAuthID       *string    `gorm:""` // OAuth provider user ID
}

type SignUpInput struct {
	Username        string `validate:"required,min=3,max=20" form:"username"`
	Email           string `validate:"required,email" form:"email"`
	Password        string `validate:"required,min=8" form:"password"`
	PasswordConfirm string `validate:"required,min=8" form:"confirm_password"`
}
