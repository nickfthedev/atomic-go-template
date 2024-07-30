package handler

import (
	"my-go-template/cmd/web/components"
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
	templ.Handler(components.SuccessResponse(components.SuccessResponseData{
		Message: "Signup successful. A verification email has been sent to your email address. Please verify your email address to continue.",
		ActionButton: &components.ActionButton{
			Label: "Login",
			Url:   "/auth/login",
		},
	})).ServeHTTP(w, r)
}

func (h *Handler) HandleForgetPassword(w http.ResponseWriter, r *http.Request) {
	templ.Handler(components.SuccessResponse(components.SuccessResponseData{
		Message: "A password reset email has been sent to your email address. Please check your email to reset your password.",
	})).ServeHTTP(w, r)
}
