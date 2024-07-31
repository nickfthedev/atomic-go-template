package middleware

import (
	"my-go-template/internal/database"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

type Middleware struct {
	db          database.Service
	validate    *validator.Validate
	formDecoder *form.Decoder
}

func NewMiddleware(db database.Service, validate *validator.Validate, formDecoder *form.Decoder) *Middleware {
	return &Middleware{db: db, validate: validate, formDecoder: formDecoder}
}
