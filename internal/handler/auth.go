package handler

import (
	"fmt"
	"my-go-template/cmd/web/auth"
	"my-go-template/cmd/web/components"
	"my-go-template/internal/model"
	"my-go-template/internal/utils"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// HandleSignup handles the signup form submission, validates the input and creates a new user
func (h *Handler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	// Declare the input struct
	var input model.SignUpInput

	// Parse and bind the form data to the input struct
	if err := utils.ParseAndBindForm(r, &input, h.formDecoder); err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Error processing form data: " + err.Error()},
		}))).ServeHTTP(w, r)
		return
	}

	// Validate the input
	if err := h.validate.Struct(input); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var messages []string
		for _, validationError := range validationErrors {
			messages = append(messages, utils.MsgForTag(validationError))
		}
		// Handle validation errors
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: messages,
		}))).ServeHTTP(w, r)
		return
	}

	// Check if password and confirm password match
	if input.Password != input.PasswordConfirm {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Passwords do not match"},
		}))).ServeHTTP(w, r)
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Error hashing password: " + err.Error()},
		}))).ServeHTTP(w, r)
		return
	}

	// Save user to database
	verifyMailToken := uuid.New().String()
	user := model.User{
		Username:          input.Username,
		Email:             input.Email,
		VerifyMailAddress: &input.Email,
		VerifyMailToken:   &verifyMailToken,
		Password:          &hashedPassword,
	}
	if err := h.db.GetDB().Create(&user).Error; err != nil {
		// Check if it's a unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
				Messages: []string{"A user with this email or username already exists"},
			}))).ServeHTTP(w, r)
		} else {
			addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
				Messages: []string{"Error creating user: " + err.Error()},
			}))).ServeHTTP(w, r)
		}
		return
	}
	if h.config.Auth.EnableVerifyEmail {
		// Send verification email
		err := h.mail.Send(user.Email, h.config.App.Name+" - Verify your email address", "Thank you for signing up. Please click the link below to verify your email address: "+h.config.App.Url+"/auth/verify-email?token="+verifyMailToken)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	// Return a success response
	if h.config.Auth.EnableVerifyEmail {
		addSuccessHeaderHandler(templ.Handler(components.SuccessResponse(components.SuccessResponseData{
			Message: "Signup successful. A verification email has been sent to your email address. Please verify your email address to continue.",
			ActionButton: &components.ActionButton{
				Label: "Login",
				Url:   "/auth/login",
			},
		}))).ServeHTTP(w, r)
	} else {
		addSuccessHeaderHandler(templ.Handler(components.SuccessResponse(components.SuccessResponseData{
			Message: "Signup successful. You can now login.",
		}))).ServeHTTP(w, r)
	}
}

func (h *Handler) HandleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Get Token from URL
	token := r.URL.Query().Get("token")

	// Verify Token
	user := model.User{}
	if err := h.db.GetDB().First(&user, "verify_mail_token = ?", token).Error; err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBannerFullPage(components.ErrorBannerData{
			Messages: []string{"Invalid verification token"},
		}))).ServeHTTP(w, r)
		return
	}

	// If found set verifiedAt to the current Date
	if user.VerifiedAt == nil {
		h.db.GetDB().Model(&user).Updates(map[string]interface{}{
			"email":             *user.VerifyMailAddress,
			"verified_at":       time.Now(),
			"verify_mail_token": nil,
		})
	}

	templ.Handler(components.SuccessResponseFullPage(components.SuccessResponseData{
		Message:      "Email verified successfully. You will be redirected to the login page in 2 seconds.",
		RedirectUrl:  &[]string{"/auth/login"}[0],
		RedirectTime: &[]int{2}[0],
	})).ServeHTTP(w, r)
}

func (h *Handler) HandleForgetPassword(w http.ResponseWriter, r *http.Request) {
	var input model.ForgotPasswordInput
	if err := utils.ParseAndBindForm(r, &input, h.formDecoder); err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Error processing form data: " + err.Error()},
		}))).ServeHTTP(w, r)
		return
	}

	// Validate the input
	if err := h.validate.Struct(input); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var messages []string
		for _, validationError := range validationErrors {
			messages = append(messages, utils.MsgForTag(validationError))
		}
		// Handle validation errors
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: messages,
		}))).ServeHTTP(w, r)
		return
	}

	// Check if the user exists, we dont want to handle the err because the user should not know if the email is valid or not
	user := model.User{}
	h.db.GetDB().First(&user, "email = ?", input.Email)

	if user.ID != uuid.Nil {
		// Generate a random password reset token
		token := uuid.New().String()
		// Update the user with the password reset token
		user.PasswordResetToken = &token
		now := time.Now()
		user.PasswordResetRequestedAt = &now
		h.db.GetDB().Save(&user)

		// Send verification email
		err := h.mail.Send(user.Email,
			h.config.App.Name+" - Reset your password",
			"Please click the link below to reset your password: "+h.config.App.Url+"/auth/reset-password?token="+*user.PasswordResetToken,
		)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} // END IF USER EXISTS

	addSuccessHeaderHandler(templ.Handler(components.SuccessResponse(components.SuccessResponseData{
		Message: "A password reset email has been sent to your email address. Please check your email to reset your password. The link is valid for 24 hours",
	}))).ServeHTTP(w, r)
}

func (h *Handler) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	// Get the token from the URL
	token := r.URL.Query().Get("token")

	// Find the user with the token
	user := model.User{}
	err := h.db.GetDB().First(&user, "password_reset_token = ?", token).Error
	if err != nil {
		templ.Handler(components.ErrorBannerFullPage(components.ErrorBannerData{
			Messages: []string{"Invalid password reset token. Maybe expired?"},
		})).ServeHTTP(w, r)
		return
	}

	// Check if link is valid
	if user.PasswordResetRequestedAt == nil || user.PasswordResetRequestedAt.Add(24*time.Hour).Before(time.Now()) {
		templ.Handler(components.ErrorBannerFullPage(components.ErrorBannerData{
			Messages: []string{"Invalid password reset token. Maybe expired?"},
		})).ServeHTTP(w, r)
		return
	}

	templ.Handler(auth.ResetPasswordForm(r)).ServeHTTP(w, r)
}

func (h *Handler) HandleResetPasswordSubmit(w http.ResponseWriter, r *http.Request) {

	var input model.ResetPasswordInput
	if err := utils.ParseAndBindForm(r, &input, h.formDecoder); err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Error processing form data: " + err.Error()},
		}))).ServeHTTP(w, r)
		return
	}

	// Validate the input
	if err := h.validate.Struct(input); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var messages []string
		for _, validationError := range validationErrors {
			messages = append(messages, utils.MsgForTag(validationError))
		}
		// Handle validation errors
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: messages,
		}))).ServeHTTP(w, r)
		return
	}

	if input.Password != input.PasswordConfirm {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Passwords do not match"},
		}))).ServeHTTP(w, r)
		return
	}

	// Find the user with the token
	user := model.User{}
	err := h.db.GetDB().First(&user, "password_reset_token = ?", input.Token).Error
	if err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Invalid password reset token. Maybe expired?"},
		}))).ServeHTTP(w, r)
		return
	}

	// Check if link is valid
	if user.PasswordResetRequestedAt == nil || user.PasswordResetRequestedAt.Add(24*time.Hour).Before(time.Now()) {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Invalid password reset token. Maybe expired?"},
		}))).ServeHTTP(w, r)
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Error hashing password: " + err.Error()},
		}))).ServeHTTP(w, r)
		return
	}

	// Update the user
	user.Password = &hashedPassword
	user.PasswordResetRequestedAt = nil
	user.PasswordResetToken = nil
	h.db.GetDB().Save(&user)

	addSuccessHeaderHandler(templ.Handler(components.SuccessResponse(components.SuccessResponseData{
		Message:      "Password reset successful. You will be redirected to the login page in 2 seconds.",
		RedirectUrl:  &[]string{"/auth/login"}[0],
		RedirectTime: &[]int{2}[0],
	}))).ServeHTTP(w, r)

}
