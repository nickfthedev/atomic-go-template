package handler

import (
	"my-go-template/cmd/web/components"
	"net/http"

	"github.com/a-h/templ"
)

func (h *Handler) Theme(w http.ResponseWriter, r *http.Request) {
	theme := r.FormValue("theme")

	if theme == "system" {
		http.SetCookie(w, &http.Cookie{
			Name:  "theme",
			Path:  "/",
			Value: "",
		})
	}
	if theme == "light" {
		http.SetCookie(w, &http.Cookie{
			Name:  "theme",
			Path:  "/",
			Value: "light",
		})
	}
	if theme == "dark" {
		http.SetCookie(w, &http.Cookie{
			Name:  "theme",
			Path:  "/",
			Value: "dark",
		})
	}

	templ.Handler(components.RedirectResponse(components.RedirectResponseData{
		RedirectUrl:  "", // Leave empty for redirect to current route
		RedirectTime: "0",
	})).ServeHTTP(w, r)
}
