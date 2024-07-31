package server

import (
	"net/http"

	"my-go-template/cmd/web"
	"my-go-template/cmd/web/auth"
	"my-go-template/cmd/web/embed"
	"my-go-template/internal/handler"

	mw "my-go-template/internal/middleware"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	// Logging Middleware
	r.Use(middleware.Logger)

	// Create a new handler instance
	h := handler.NewHandler(s.db, s.validate, s.formDecoder)
	// Create a new middleware instance for own middlewares
	m := mw.NewMiddleware(s.db, s.validate, s.formDecoder)

	// Checks for the JWT token in the cookie and sets the user data into the context
	r.Use(m.JWTMiddleware)

	// Serve static files
	fileServer := http.FileServer(http.FS(embed.Files))
	r.Handle("/assets/*", fileServer)

	// API Test Endpoint
	r.Get("/api", h.HelloWorldHandler)
	// Health Check
	r.Get("/health", h.HealthHandler)
	// Hello
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(web.HelloForm(r)).ServeHTTP(w, r)
	})
	r.Post("/hello", h.HelloWebHandler)

	// Theme
	r.Post("/theme", h.Theme)
	// Auth Group goes here
	r.Route("/auth", func(r chi.Router) {
		r.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
			templ.Handler(auth.SignupForm(r)).ServeHTTP(w, r)
		})
		r.Get("/verify-email", h.HandleVerifyEmail)
		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			templ.Handler(auth.LoginForm(r)).ServeHTTP(w, r)
		})
		r.Get("/logout", h.HandleLogout)
		r.Get("/forget-password", func(w http.ResponseWriter, r *http.Request) {
			templ.Handler(auth.ForgetPasswordForm(r)).ServeHTTP(w, r)
		})
		r.Post("/login", h.HandleLogin)
		r.Post("/signup", h.HandleSignup)
		r.Post("/forget-password", h.HandleForgetPassword)
		r.Get("/reset-password", h.HandleResetPassword)
		r.Post("/reset-password", h.HandleResetPasswordSubmit)
	})
	r.Get("/profile/edit", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(auth.EditProfilePage(r)).ServeHTTP(w, r)
	})
	r.Post("/profile/edit", h.HandleEditProfile)
	return r
}
