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

		_, ok := r.Context().Value(UserKey).(model.User)
		if !ok {
			// Unset the JWT cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "jwt",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				MaxAge:   -1,
			})
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}

		next(w, r)
	}
}
