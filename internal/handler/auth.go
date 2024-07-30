package handler

import (
	"fmt"
	"my-go-template/cmd/web/components"
	"my-go-template/internal/model"
	"net/http"

	"github.com/a-h/templ"
)

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	templ.Handler(components.SuccessResponse(components.SuccessResponseData{
		Message:      "Login successful",
		RedirectUrl:  &[]string{"/"}[0],
		RedirectTime: &[]int{2}[0],
	})).ServeHTTP(w, r)
}

func (h *Handler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var input model.SignUpInput
	if err := r.ParseForm(); err != nil {
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Message: "Error parsing form data: " + err.Error(),
		}))).ServeHTTP(w, r)
		return
	}
	input.Username = r.FormValue("username")
	input.Email = r.FormValue("email")
	input.Password = r.FormValue("password")
	input.PasswordConfirm = r.FormValue("confirm_password")
	fmt.Println(input)

	// TODO: Human readable error messages

	// Validate the input
	if err := h.validate.Struct(input); err != nil {
		// Handle validation errors
		addErrorHeaderHandler(templ.Handler(components.ErrorBanner(components.ErrorBannerData{
			Message: "Invalid input: " + err.Error(),
		}))).ServeHTTP(w, r)
		return
	}

	// Return a success response
	addSuccessHeaderHandler(templ.Handler(components.SuccessResponse(components.SuccessResponseData{
		Message: "Signup successful. A verification email has been sent to your email address. Please verify your email address to continue.",
		ActionButton: &components.ActionButton{
			Label: "Login",
			Url:   "/auth/login",
		},
	}))).ServeHTTP(w, r)
}

func (h *Handler) HandleForgetPassword(w http.ResponseWriter, r *http.Request) {
	templ.Handler(components.SuccessResponse(components.SuccessResponseData{
		Message: "A password reset email has been sent to your email address. Please check your email to reset your password.",
	})).ServeHTTP(w, r)
}
