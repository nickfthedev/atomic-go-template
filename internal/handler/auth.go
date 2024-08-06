package handler

import (
	"fmt"
	"io"
	"my-go-template/cmd/web/auth"
	"my-go-template/cmd/web/components"
	mw "my-go-template/internal/middleware"
	"my-go-template/internal/model"
	"my-go-template/internal/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

	if user.VerifiedAt == nil && h.config.Auth.EnableVerifyEmail {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Please verify your email address before logging in"},
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

func (h *Handler) HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	// Parse Multipart Form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Error processing form data: " + err.Error()},
		})).ServeHTTP(w, r)
		return
	}

	var avatarPath string
	// Check if avatar is set
	if r.MultipartForm.File["avatar"] != nil && h.config.Auth.EnableAvatar {
		avatarFile := r.MultipartForm.File["avatar"][0]

		// Open File
		avatar, err := avatarFile.Open()
		if err != nil {
			fmt.Println("Error opening avatar file:", err)
			templ.Handler(components.ErrorBanner(components.ErrorBannerData{
				Messages: []string{"Error processing form data: " + err.Error()},
			})).ServeHTTP(w, r)
			return
		}
		defer avatar.Close()

		// Read the first 512 bytes to detect content type
		buffer := make([]byte, 512)
		_, err = avatar.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading avatar file:", err)
			templ.Handler(components.ErrorBanner(components.ErrorBannerData{
				Messages: []string{"Error processing form data: " + err.Error()},
			})).ServeHTTP(w, r)
			return
		}

		// Detect content type
		contentType := http.DetectContentType(buffer)

		// Check if the content type is an allowed image format
		if contentType != "image/jpeg" && contentType != "image/png" {
			templ.Handler(components.ErrorBanner(components.ErrorBannerData{
				Messages: []string{"Invalid file format. Please upload a PNG or JPEG image."},
			})).ServeHTTP(w, r)
			return
		}

		// Reset the file pointer to the beginning
		avatar.Seek(0, 0)

		// Ensure the directory exists
		if err := os.MkdirAll("public/avatars", os.ModePerm); err != nil {
			fmt.Println("Error creating directory:", err)
			templ.Handler(components.ErrorBanner(components.ErrorBannerData{
				Messages: []string{"Error processing form data: " + err.Error()},
			})).ServeHTTP(w, r)
			return
		}

		// Generate a random filename
		avatarPath = uuid.New().String() + "-" + time.Now().Format("060102150405") + "." + strings.Split(avatarFile.Filename, ".")[1]

		// Save Avatar to public folder
		file, err := os.Create("public/avatars/" + avatarPath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			templ.Handler(components.ErrorBanner(components.ErrorBannerData{
				Messages: []string{"Error processing form data: " + err.Error()},
			})).ServeHTTP(w, r)
			return
		}
		defer file.Close()

		// Copy the file content
		if _, err := io.Copy(file, avatar); err != nil {
			fmt.Println("Error copying file content:", err)
			templ.Handler(components.ErrorBanner(components.ErrorBannerData{
				Messages: []string{"Error processing form data: " + err.Error()},
			})).ServeHTTP(w, r)
			return
		}
	}

	// Declare the input struct
	var input model.EditProfileInput
	// Parse and bind the form data to the input struct
	if err := utils.ParseAndBindForm(r, &input, h.formDecoder); err != nil {
		templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Error processing form data: " + err.Error()},
		})).ServeHTTP(w, r)
		return
	}
	// Handle empty password fields
	if input.Password != nil && *input.Password == "" {
		input.Password = nil
	}
	// Validate the input
	if err := h.validate.Struct(input); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var messages []string
		for _, validationError := range validationErrors {
			messages = append(messages, utils.MsgForTag(validationError))
		}
		// Handle validation errors
		templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: messages,
		})).ServeHTTP(w, r)
		return
	}

	// Get user from database
	user := r.Context().Value(mw.UserKey).(model.User)
	user.Password = nil
	// Update User Name and Avatar Path, nil values will be skipped on saving
	updateFields := map[string]interface{}{
		"username": input.Username,
	}

	// Only update avatar_url if a new avatar is uploaded
	if avatarPath != "" {
		updateFields["avatar_url"] = avatarPath
	}

	// Check if email has changed and resend verification email, if enabled
	var verifyMailToken string
	if user.Email != input.Email {
		if h.config.Auth.EnableVerifyEmail {
			verifyMailToken = uuid.New().String()
			updateFields["verify_mail_token"] = &verifyMailToken
			updateFields["verify_mail_address"] = &input.Email
			updateFields["verified_at"] = nil
			// Verification Mail will be sent after the user is updated successfully
		} else {
			updateFields["email"] = input.Email
		}
	}

	// Check if input password is not nil
	if input.Password != nil && input.PasswordConfirm != nil {
		// Check if password and confirm password match
		if *input.Password != *input.PasswordConfirm {
			templ.Handler(components.ErrorBanner(components.ErrorBannerData{
				Messages: []string{"Passwords do not match"},
			})).ServeHTTP(w, r)
			return
		}
		hashedPassword, err := utils.HashPassword(*input.Password)
		if err != nil {
			templ.Handler(components.ErrorBanner(components.ErrorBannerData{
				Messages: []string{"Error hashing password: " + err.Error()},
			})).ServeHTTP(w, r)
			return
		}
		user.Password = &hashedPassword
	}

	// Conditionally add the password field if it has been set
	if user.Password != nil {
		updateFields["password"] = *user.Password
	}
	// Save user to database
	err := h.db.GetDB().Model(&user).Updates(updateFields).Error
	if err != nil {
		templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Messages: []string{"Error updating user: " + err.Error()},
		})).ServeHTTP(w, r)
		return
	}
	if user.Email != input.Email && h.config.Auth.EnableVerifyEmail {
		// Send verification email
		err := h.mail.Send(user.Email,
			h.config.App.Name+" - Verify your new email address",
			"Please click the link below to verify your new email address: "+h.config.App.Url+"/auth/verify-email?token="+verifyMailToken,
		)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	// Return a success response
	templ.Handler(components.SuccessResponse(components.SuccessResponseData{
		Message: "Profile updated successfully. ",
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
