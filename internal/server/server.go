package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"

	"atomic-go-template/internal/config"
	"atomic-go-template/internal/database"
	"atomic-go-template/internal/mail"
)

type Server struct {
	port int
	// The database instance
	db database.Service
	// The validator instance
	validate *validator.Validate
	// The form decoder instance
	formDecoder *form.Decoder
	// The config instance
	config *config.Config
	// The mail service instance
	mail mail.Service
}

func NewServer() *http.Server {
	// Get Port from Environment Variables
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Create Config
	// Can be used in middleware and handler as well
	// If you dont want to change values its safe to remove them here.
	// Default values are set in config.go
	config := config.New(&config.Config{
		Server: config.Server{
			Port: port,
		},
		Database: config.Database{
			Enabled: true,
			Type:    config.DatabaseTypeSQLite,
		},
		Theme: config.Theme{
			StandardTheme:       "",
			EnableThemeSwitcher: true,
			EnableSidebar:       true,
		},
		Auth: config.Auth{
			EnableAuth:          true,
			EnableRegistration:  true,
			EnableLogin:         true,
			EnableAvatar:        true,
			EnableResetPassword: true,
			EnableVerifyEmail:   true,
		},
		Mail: config.Mail{
			EnableMail:   true,
			MailProvider: config.MailProviderConsole,
		},
	})

	// Create database service
	db := database.New(config.Database)
	// Automigrate
	database.MigrateUserSchema(db.GetDB())

	// Mail Service
	var mailService mail.Service
	var err error
	if config.Mail.EnableMail {
		mailService, err = mail.NewMailProvider(config.Mail)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create server struct
	NewServer := &Server{
		port:        config.Server.Port,
		db:          db,
		validate:    validator.New(validator.WithRequiredStructEnabled()),
		formDecoder: form.NewDecoder(),
		config:      config,
		mail:        mailService,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	fmt.Printf("Server is running on port: %d", NewServer.port)
	return server
}
