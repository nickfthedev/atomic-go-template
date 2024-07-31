package handler

import (
	"fmt"
	"my-go-template/cmd/web/auth"
	"my-go-template/cmd/web/components"
	"my-go-template/internal/model"
	"my-go-template/internal/utils"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/resend/resend-go/v2"
)

// HandleLogin handles the login form submission, validates the input and redirects to the home page after successful login
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Declare the input struct
	var input model.LoginInput

	// Parse and bind the form data to the input struc
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

	// Find user in database
	user := model.User{}
	if err := h.db.GetDB().First(&user, "email = ?", input.Email).Error; err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Invalid email or password"},
		}))).ServeHTTP(w, r)
		return
	}

	if err := utils.CheckPasswordHash(*user.Password, input.Password); err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Invalid email or password"},
		}))).ServeHTTP(w, r)
		return
	}

	// Set Cookie
	if err := utils.CreateJWTCookie(w, user.ID.String()); err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Error creating JWT cookie: " + err.Error()},
		}))).ServeHTTP(w, r)
		return
	}

	// Show Success Message and send redirect to client
	addSuccessHeaderHandler(templ.Handler(components.SuccessResponse(components.SuccessResponseData{
		Message:      "Login successful",
		RedirectUrl:  &[]string{"/"}[0],
		RedirectTime: &[]int{2}[0],
	}))).ServeHTTP(w, r)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	utils.DeleteJWTCookie(w)
	templ.Handler(components.SuccessResponseFullPage(components.SuccessResponseData{
		Message:      "Logout successful. You will be redirected to the home page in 2 seconds.",
		RedirectUrl:  &[]string{"/"}[0],
		RedirectTime: &[]int{2}[0],
	})).ServeHTTP(w, r)
}

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
	user := model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: &hashedPassword,
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

	// Send verification email
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", os.Getenv("APP_NAME"), os.Getenv("RESEND_FROM_EMAIL")),
		To:      []string{user.Email},
		Html:    "Thank you for signing up. Please click the link below to verify your email address: " + os.Getenv("APP_URL") + "/auth/verify-email?token=" + user.ID.String(),
		Subject: fmt.Sprintf("%s - Verify your email address", os.Getenv("APP_NAME")),
		// Cc:      []string{"cc@example.com"},
		// Bcc:     []string{"bcc@example.com"},
		// ReplyTo: "replyto@example.com",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Verification email sent with ID:", sent.Id)

	// Return a success response
	addSuccessHeaderHandler(templ.Handler(components.SuccessResponse(components.SuccessResponseData{
		Message: "Signup successful. A verification email has been sent to your email address. Please verify your email address to continue.",
		ActionButton: &components.ActionButton{
			Label: "Login",
			Url:   "/auth/login",
		},
	}))).ServeHTTP(w, r)
}

func (h *Handler) HandleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Get Token from URL
	token := r.URL.Query().Get("token")

	// Verify Token
	user := model.User{}
	if err := h.db.GetDB().First(&user, "id = ?", token).Error; err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBannerFullPage(components.ErrorBannerData{
			Messages: []string{"Invalid verification token"},
		}))).ServeHTTP(w, r)
		return
	}

	// If found set verifiedAt to the current Date
	if user.VerifiedAt == nil {
		h.db.GetDB().Model(&user).Update("verified_at", time.Now())
		h.db.GetDB().Save(&user)
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
		client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
		params := &resend.SendEmailRequest{
			From:    fmt.Sprintf("%s <%s>", os.Getenv("APP_NAME"), os.Getenv("RESEND_FROM_EMAIL")),
			To:      []string{user.Email},
			Html:    "Please click the link below to reset your password: " + os.Getenv("APP_URL") + "/auth/reset-password?token=" + *user.PasswordResetToken,
			Subject: fmt.Sprintf("%s - Reset your password", os.Getenv("APP_NAME")),
			// Cc:      []string{"cc@example.com"},
			// Bcc:     []string{"bcc@example.com"},
			// ReplyTo: "replyto@example.com",
		}

		sent, err := client.Emails.Send(params)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Verification email sent with ID:", sent.Id)
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

	// Get the token from the URL
	token := r.URL.Query().Get("token") // TODO: Token not sent with HTMX
	// Find the user with the token
	user := model.User{}
	err := h.db.GetDB().First(&user, "password_reset_token = ?", token).Error
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
