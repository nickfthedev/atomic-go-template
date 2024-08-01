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
	h := handler.NewHandler(s.db, s.validate, s.formDecoder, s.config)
	// Create a new middleware instance for own middlewares
	m := mw.NewMiddleware(s.db, s.validate, s.formDecoder, s.config)

	// Add Config to Context
	r.Use(m.ConfigMiddleware)

	// Checks for the JWT token in the cookie and sets the user data into the context
	r.Use(m.JWTMiddleware)

	// Serve static files
	fileServer := http.FileServer(http.FS(embed.Files))
	r.Handle("/assets/*", fileServer)

	// Public Folder // TODO: Figure out how to preserve this folder in a Dockerfile
	publicFileServer := http.FileServer(http.Dir("cmd/web/public"))
	r.Handle("/public/*", http.StripPrefix("/public", publicFileServer))

	// API Test Endpoint
	r.Get("/api", h.HelloWorldHandler)

	// Health Check
	r.Get("/health", h.HealthHandler)

	// Home
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(web.HelloForm(r)).ServeHTTP(w, r)
	})
	r.Post("/hello", h.HelloWebHandler) // Test on Homepage
	// Sample Protected Page
	r.Get("/protected", m.IsLoggedIn(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Protected. Access granted"))
	}))

	// Theme
	if s.config.Theme.EnableThemeSwitcher {
		r.Post("/theme", h.Theme)
	}
	// Auth Group goes here
	if s.config.Auth.EnableAuth {
		r.Route("/auth", func(r chi.Router) {
			// Signup Routes
			if s.config.Auth.EnableRegistration {
				r.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
					templ.Handler(auth.SignupForm(r)).ServeHTTP(w, r)
				})
			}
			// Login Routes
			if s.config.Auth.EnableLogin {
				r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
					templ.Handler(auth.LoginForm(r)).ServeHTTP(w, r)
				})
				r.Post("/login", h.HandleLogin)
				r.Get("/logout", h.HandleLogout)
			}
			// Reset Password Routes
			if s.config.Auth.EnableResetPassword {
				r.Get("/forget-password", func(w http.ResponseWriter, r *http.Request) {
					templ.Handler(auth.ForgetPasswordForm(r)).ServeHTTP(w, r)
				})
				r.Post("/forget-password", h.HandleForgetPassword)
				r.Get("/reset-password", h.HandleResetPassword)
				r.Post("/reset-password", h.HandleResetPasswordSubmit)
			}
			// Verify Email Routes
			if s.config.Auth.EnableVerifyEmail {
				r.Get("/verify-email", h.HandleVerifyEmail)
			}
		}) // End of Auth Group

		// Profile Routes
		r.Get("/profile/edit", m.IsLoggedIn(func(w http.ResponseWriter, r *http.Request) {
			templ.Handler(auth.EditProfilePage(r)).ServeHTTP(w, r)
		}))
		r.Post("/profile/edit", m.IsLoggedIn(h.HandleEditProfile))
	} // End of Auth Feature Routes
	return r
}
