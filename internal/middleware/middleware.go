package middleware

import (
	"atomic-go-template/internal/config"
	"atomic-go-template/internal/database"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

type ContextKey string

type Middleware struct {
	db          database.Service
	validate    *validator.Validate
	formDecoder *form.Decoder
	config      *config.Config
}

func NewMiddleware(db database.Service, validate *validator.Validate, formDecoder *form.Decoder, config *config.Config) *Middleware {
	return &Middleware{db: db, validate: validate, formDecoder: formDecoder, config: config}
}
