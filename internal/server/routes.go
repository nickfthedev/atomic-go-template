package server

import (
	"net/http"
	"os"
	"strings"

	"my-go-template/cmd/web/auth"
	"my-go-template/internal/handler"
	"my-go-template/web/embed"
	"my-go-template/web/routes"
	"my-go-template/web/routes/auth/login"
	"my-go-template/web/routes/auth/logout"
	"my-go-template/web/routes/protected"
	"my-go-template/web/routes/user/profile"

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
	h := handler.NewHandler(s.db, s.validate, s.formDecoder, s.config, s.mail)
	// Create a new middleware instance for own middlewares
	m := mw.NewMiddleware(s.db, s.validate, s.formDecoder, s.config)

	// Add Config to Context
	r.Use(m.ConfigMiddleware)

	// Checks for the JWT token in the cookie and sets the user data into the context
	r.Use(m.JWTMiddleware)

	// Serve static files without directory listing
	fileServer := http.FileServer(NoListingFileSystem{http.FS(embed.Files)})
	r.Handle("/assets/*", fileServer)

	// Public Folder without directory listing
	publicFileServer := http.FileServer(NoListingFileSystem{http.Dir("public")})
	r.Handle("/public/*", http.StripPrefix("/public", publicFileServer))

	// API Test Endpoint
	r.Get("/api", h.HelloWorldHandler)

	// Health Check
	r.Get("/health", h.HealthHandler)

	// Home
	r.Get("/", routes.GET)
	r.Post("/hello", h.HelloWebHandler) // Test on Homepage
	// Sample Protected Page

	// This route is only accessible if the user is logged in
	r.Get("/protected", m.IsLoggedIn(protected.New().GET))

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
				r.Post("/signup", h.HandleSignup)
			}
			// Login Routes
			if s.config.Auth.EnableLogin {
				r.Get("/login", login.New(s.db.GetDB(), s.config, s.validate, s.formDecoder).GET)
				r.Post("/login", login.New(s.db.GetDB(), s.config, s.validate, s.formDecoder).POST)
				r.Get("/logout", logout.New().GET)
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
		r.Get("/profile/edit", m.IsLoggedIn(profile.New(s.db.GetDB(), s.config, s.validate, s.formDecoder, s.mail).GET))
		r.Post("/profile/edit", m.IsLoggedIn(profile.New(s.db.GetDB(), s.config, s.validate, s.formDecoder, s.mail).POST))
	} // End of Auth Feature Routes
	return r
}

// NoListingFileSystem wraps http.FileSystem to disable directory listing
type NoListingFileSystem struct {
	fs http.FileSystem
}

func (nfs NoListingFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := nfs.fs.Open(index); err != nil {
			return nil, os.ErrNotExist
		}
	}

	return f, nil
}
