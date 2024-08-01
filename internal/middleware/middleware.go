package middleware

import (
	"my-go-template/internal/config"
	"my-go-template/internal/database"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

type Middleware struct {
	db          database.Service
	validate    *validator.Validate
	formDecoder *form.Decoder
	config      *config.Config
}

func NewMiddleware(db database.Service, validate *validator.Validate, formDecoder *form.Decoder, config *config.Config) *Middleware {
	return &Middleware{db: db, validate: validate, formDecoder: formDecoder, config: config}
}
