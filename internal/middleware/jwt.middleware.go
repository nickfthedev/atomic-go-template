package middleware

import (
	"context"
	"my-go-template/internal/model"
	"my-go-template/internal/utils"
	"net/http"
)

const UserIDKey ContextKey = "userid"
const UserKey ContextKey = "user"

// HTTP middleware setting a value on the request context
func (m *Middleware) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.VerifyJWTCookie(r)
		if err != nil {
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// Get User from DB
		// TODO: Clear out Passwords or use another struct
		user := model.User{}
		err = m.db.GetDB().First(&user, "id = ?", userID).Error
		if err != nil {
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			next.ServeHTTP(w, r)
			return
		}

		// Clear out password
		user.Password = nil

		ctx = context.WithValue(ctx, UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
