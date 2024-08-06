package handler

import (
	"my-go-template/internal/config"
	"my-go-template/internal/database"
	"my-go-template/internal/mail"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	db          database.Service
	validate    *validator.Validate
	formDecoder *form.Decoder
	config      *config.Config
	mail        mail.Service
}

func NewHandler(db database.Service, validate *validator.Validate, formDecoder *form.Decoder, config *config.Config, mail mail.Service) *Handler {
	return &Handler{db: db, validate: validate, formDecoder: formDecoder, config: config, mail: mail}
}

// TODO: Remove after rebuild
// We use this to re-target the errors div and swap the innerHTML, instead of the default behavior of appending to the end of the div.
// Useful for scenarios where you switch innerhtml for example at sign up.
// We dont want people to sign up multiple times so the form shuold disappear on success.
// On error we want to keep the form there so they can fix it. Thats why we need to Reswap the target to a errors div.
func addErrorHeaderHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("HX-Retarget", "#errors")
		w.Header().Add("HX-Reswap", "innerHTML")
		handler.ServeHTTP(w, r)
	})
}

// TODO: Remove after rebuild
// We use this to trigger the clearErrors event, which clears the errors div.
// This is useful for scenarios where you want to clear the errors div after a successful action.
func addSuccessHeaderHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("HX-Trigger", "clearErrors")
		handler.ServeHTTP(w, r)
	})
}
