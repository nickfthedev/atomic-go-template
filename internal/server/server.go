package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"

	"my-go-template/internal/database"
)

// TODO Config with Features Flags

type Server struct {
	port int
	// The database instance
	db database.Service
	// The validator instance
	validate *validator.Validate
	// The form decoder instance
	formDecoder *form.Decoder
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()                   // Create database service
	database.MigrateUserSchema(db.GetDB()) // Automigrate
	// Create server struct
	NewServer := &Server{
		port:        port,
		db:          db,
		validate:    validator.New(validator.WithRequiredStructEnabled()),
		formDecoder: form.NewDecoder(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
