package middleware

import (
	"context"
	"my-go-template/internal/config"
	"net/http"
)

const ConfigKey ContextKey = "config"

// HTTP middleware setting a value on the request context
func (m *Middleware) ConfigMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := context.WithValue(r.Context(), ConfigKey, m.config)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetConfigFromContext returns the config from the context
func GetConfigFromContext(r *http.Request) *config.Config {
	config, ok := r.Context().Value(ConfigKey).(*config.Config)
	if !ok {
		return nil
	}
	return config
}
