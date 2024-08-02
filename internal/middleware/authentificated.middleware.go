package middleware

import (
	"my-go-template/internal/model"
	"my-go-template/internal/utils"
	"net/http"
)

// IsLoggedIn checks if the user is authenticated, otherwise redirects to login or home if auth is disabled
func (m *Middleware) IsLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.config.Auth.EnableAuth {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		_, err := utils.VerifyJWTCookie(r)
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}

		// TODO: Check if the user object exists
		_, ok := r.Context().Value(UserKey).(model.User)
		if !ok {
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}

		next(w, r)
	}
}
