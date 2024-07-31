package utils

import (
	"net/http"

	"github.com/go-playground/form/v4"
)

func ParseAndBindForm(r *http.Request, form interface{}, formDecoder *form.Decoder) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	// Bind the form data to the input struct
	if err := formDecoder.Decode(&form, r.PostForm); err != nil {
		return err
	}
	return nil
}
