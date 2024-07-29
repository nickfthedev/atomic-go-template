package server

import (
	"net/http"

	"my-go-template/cmd/web"
	"my-go-template/cmd/web/auth"
	"my-go-template/cmd/web/embed"
	"my-go-template/internal/handler"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	// Logging Middleware
	r.Use(middleware.Logger)

	// Create a new handler instance
	h := handler.NewHandler(s.db)

	// Serve static files
	fileServer := http.FileServer(http.FS(embed.Files))
	r.Handle("/assets/*", fileServer)

	// API Test Endpoint
	r.Get("/api", h.HelloWorldHandler)
	// Health Check
	r.Get("/health", h.HealthHandler)
	// Hello
	r.Get("/", templ.Handler(web.HelloForm()).ServeHTTP)
	r.Post("/hello", h.HelloWebHandler)

	// Auth Group goes here
	r.Route("/auth", func(r chi.Router) {
		r.Get("/signup", templ.Handler(auth.SignupForm()).ServeHTTP)

		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Login"))
		})
	})

	return r
}
