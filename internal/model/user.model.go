package model

import "time"

// User represents a user in the database.
type User struct {
	BaseModel
	Username                 string     `gorm:"unique;not null"`
	Email                    string     `gorm:"unique;not null"`
	Password                 *string    `gorm:""` // Password is optional
	PasswordResetToken       *string    `gorm:""` // Password reset token is optional
	PasswordResetRequestedAt *time.Time `gorm:""` // Password reset requested at is optional
	VerifiedAt               *time.Time `gorm:""` // Verified at is optional
	VerifyMailAddress        *string    `gorm:""` // Verify mail address is optional
	VerifyMailToken          *string    `gorm:""` // Verify mail token is optional
	AvatarURL                *string    `gorm:""` // Avatar URL is optional
	OAuthProvider            *string    `gorm:""` // OAuth provider name (e.g., "google", "github")
	OAuthID                  *string    `gorm:""` // OAuth provider user ID
}

type SignUpInput struct {
	Username        string `validate:"required,min=3,max=20" form:"username"`
	Email           string `validate:"required,email" form:"email"`
	Password        string `validate:"required,min=8" form:"password"`
	PasswordConfirm string `validate:"required,min=8" form:"confirm_password"`
}

type EditProfileInput struct {
	Username        string  `validate:"required,min=3,max=20" form:"username"`
	Email           string  `validate:"required,email" form:"email"`
	Password        *string `validate:"omitempty,min=8" form:"password"`
	PasswordConfirm *string `validate:"-" form:"confirm_password"`
	AvatarURL       *string `validate:"omitempty" form:"avatar_url"`
}

type LoginInput struct {
	Email    string `validate:"required,email" form:"email"`
	Password string `validate:"required" form:"password"`
}

type ForgotPasswordInput struct {
	Email string `validate:"required,email" form:"email"`
}

type ResetPasswordInput struct {
	Password        string `validate:"required,min=8" form:"password"`
	PasswordConfirm string `validate:"required,min=8" form:"confirm_password"`
	Token           string `validate:"required" form:"token"`
}
