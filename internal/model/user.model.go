package model

// User represents a user in the database.
type User struct {
	BaseModel
	Username      string  `gorm:"unique;not null"`
	Email         string  `gorm:"unique;not null"`
	Password      *string `gorm:""` // Password is optional
	OAuthProvider *string `gorm:""` // OAuth provider name (e.g., "google", "github")
	OAuthID       *string `gorm:""` // OAuth provider user ID
}

type SignUpInput struct {
	Username        string `validate:"required,min=3,max=20"`
	Email           string `validate:"required,email"`
	Password        string `validate:"required,min=8"`
	PasswordConfirm string `validate:"required,min=8"`
}
