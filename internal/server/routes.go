package server

import (
	"net/http"
	"os"
	"strings"

	mw "atomic-go-template/internal/middleware"
	"atomic-go-template/web/components/theme"
	"atomic-go-template/web/embed"
	"atomic-go-template/web/routes"
	forget_password "atomic-go-template/web/routes/auth/forget_password"
	"atomic-go-template/web/routes/auth/login"
	"atomic-go-template/web/routes/auth/logout"
	reset_password "atomic-go-template/web/routes/auth/reset_password"
	"atomic-go-template/web/routes/auth/signup"
	verify_mail "atomic-go-template/web/routes/auth/verify-mail"
	"atomic-go-template/web/routes/health"
	"atomic-go-template/web/routes/protected"
	react_example "atomic-go-template/web/routes/react-example"
	"atomic-go-template/web/routes/user/profile"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	// Logging Middleware
	r.Use(middleware.Logger)

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

	// Health Check
	r.Get("/health", health.New(s.db, s.config).GET)

	// Home
	r.Get("/", routes.GET)

	// This route is only accessible if the user is logged in
	r.Get("/protected", m.IsLoggedIn(protected.New().GET))

	r.Get("/react-example", react_example.New(s.db.GetDB(), s.config, s.validate, s.formDecoder, s.mail).GET)

	// Theme
	if s.config.Theme.EnableThemeSwitcher {
		r.Post("/theme", theme.New().POST)
	}
	// Auth Group goes here
	if s.config.Auth.EnableAuth {
		r.Route("/auth", func(r chi.Router) {
			// Signup Routes
			if s.config.Auth.EnableRegistration {
				r.Get("/signup", signup.New(s.db.GetDB(), s.config, s.validate, s.formDecoder, s.mail).GET)
				r.Post("/signup", signup.New(s.db.GetDB(), s.config, s.validate, s.formDecoder, s.mail).POST)
			}
			// Login Routes
			if s.config.Auth.EnableLogin {
				r.Get("/login", login.New(s.db.GetDB(), s.config, s.validate, s.formDecoder).GET)
				r.Post("/login", login.New(s.db.GetDB(), s.config, s.validate, s.formDecoder).POST)
				r.Get("/logout", logout.New().GET)
			}
			// Reset Password Routes
			if s.config.Auth.EnableResetPassword {
				r.Get("/forget-password", forget_password.New(s.db.GetDB(), s.config, s.validate, s.formDecoder, s.mail).GET)
				r.Post("/forget-password", forget_password.New(s.db.GetDB(), s.config, s.validate, s.formDecoder, s.mail).POST)
				r.Get("/reset-password", reset_password.New(s.db.GetDB(), s.config, s.validate, s.formDecoder, s.mail).GET)
				r.Post("/reset-password", reset_password.New(s.db.GetDB(), s.config, s.validate, s.formDecoder, s.mail).POST)
			}
			// Verify Email Routes
			if s.config.Auth.EnableVerifyEmail {
				r.Get("/verify-email", verify_mail.New(s.db.GetDB()).GET)
			}
		}) // End of Auth Group

		// Profile Routes
		r.Get("/user/profile", m.IsLoggedIn(profile.New(s.db.GetDB(), s.config, s.validate, s.formDecoder, s.mail).GET))
		r.Post("/user/profile", m.IsLoggedIn(profile.New(s.db.GetDB(), s.config, s.validate, s.formDecoder, s.mail).POST))
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
