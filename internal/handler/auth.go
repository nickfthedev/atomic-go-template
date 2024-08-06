package handler

import (
	"my-go-template/cmd/web/auth"
	"my-go-template/cmd/web/components"
	"my-go-template/internal/model"
	"my-go-template/internal/utils"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
)

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
