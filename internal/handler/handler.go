package handler

import (
	"my-go-template/internal/database"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	db       database.Service
	validate *validator.Validate
}

func NewHandler(db database.Service, validate *validator.Validate) *Handler {
	return &Handler{db: db, validate: validate}
}

// We use this to re-target the errors div and swap the innerHTML, instead of the default behavior of appending to the end of the div.
func addErrorHeaderHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("HX-Retarget", "#errors")
		w.Header().Add("HX-Reswap", "innerHTML")
		handler.ServeHTTP(w, r)
	})
}

// We use this to trigger the clearErrors event, which clears the errors div.
func addSuccessHeaderHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("HX-Trigger", "clearErrors")
		handler.ServeHTTP(w, r)
	})
}
